package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"API/storage/cache"
	"API/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// RateLimiter 基础限流中间件
func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := redis_rate.NewLimiter(cache.RedisClient)
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/v1/health" {
			c.Next()
			return
		}
		res, err := limiter.Allow(
			c.Request.Context(),
			"global_rate_limit",
			redis_rate.Limit{
				Rate:   maxRequests,
				Period: window,
				Burst:  maxRequests * 2,
			},
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if res.Allowed == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "请求过于频繁，请稍后再试",
				"retry_after": res.RetryAfter,
			})
			return
		}
		c.Next()
	}
}

// RateLimiterConfig 限流器配置
type RateLimiterConfig struct {
	Enabled      bool
	MaxRequests  int
	Window       time.Duration
	Burst        int
	ExcludePaths []string
}

// LoadRateLimiterConfig 从配置文件加载限流配置
func LoadRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Enabled:      viper.GetBool("ratelimit.enabled"),
		MaxRequests:  viper.GetInt("ratelimit.max_requests"),
		Window:       viper.GetDuration("ratelimit.window") * time.Second,
		Burst:        viper.GetInt("ratelimit.burst"),
		ExcludePaths: viper.GetStringSlice("ratelimit.exclude_paths"),
	}
}

// EnhancedRateLimiter 增强版限流中间件
func EnhancedRateLimiter(logger *zap.Logger) gin.HandlerFunc {
	limiter := redis_rate.NewLimiter(cache.RedisClient)
	config := LoadRateLimiterConfig()

	// 如果未启用限流，返回空中间件
	if !config.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 确保配置有效
	if config.MaxRequests <= 0 {
		config.MaxRequests = 100 // 默认值
	}
	if config.Window <= 0 {
		config.Window = time.Minute // 默认值
	}
	if config.Burst <= 0 {
		config.Burst = config.MaxRequests * 2 // 默认值
	}

	logger.Info("限流器已启用",
		zap.Int("max_requests", config.MaxRequests),
		zap.Duration("window", config.Window),
		zap.Int("burst", config.Burst),
	)

	return func(c *gin.Context) {
		// 检查是否为排除路径
		for _, path := range config.ExcludePaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// 获取客户端IP作为限流键
		clientIP := c.ClientIP()
		key := "rate_limit:" + clientIP

		// 应用限流
		res, err := limiter.Allow(
			c.Request.Context(),
			key,
			redis_rate.Limit{
				Rate:   config.MaxRequests,
				Period: config.Window,
				Burst:  config.Burst,
			},
		)

		// 处理Redis错误
		if err != nil {
			logger.Error("限流器Redis错误", zap.Error(err))
			// 降级处理：允许请求通过
			c.Next()
			return
		}

		// 如果超出限制
		if res.Allowed == 0 {
			logger.Warn("请求被限流",
				zap.String("ip", clientIP),
				zap.String("path", c.Request.URL.Path),
				zap.Duration("retry_after", res.RetryAfter),
			)

			// 设置重试头
			c.Header("Retry-After", res.RetryAfter.String())
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.MaxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(res.RetryAfter).Unix()))

			// 使用增强的错误响应
			utils.RespondWithJSON(
				c,
				http.StatusTooManyRequests,
				429,
				"请求过于频繁，请稍后再试",
				gin.H{"retry_after": res.RetryAfter.Seconds()},
			)
			c.Abort()
			return
		}

		// 设置限流头信息
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.MaxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))

		c.Next()
	}
}

// IPBasedRateLimiter IP地址限流中间件
func IPBasedRateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := redis_rate.NewLimiter(cache.RedisClient)
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "ip_rate_limit:" + clientIP

		res, err := limiter.Allow(
			context.Background(),
			key,
			redis_rate.Limit{
				Rate:   maxRequests,
				Period: window,
				Burst:  maxRequests,
			},
		)

		if err != nil || res.Allowed > 0 {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error":       "IP请求过于频繁，请稍后再试",
			"retry_after": res.RetryAfter.Seconds(),
		})
	}
}
