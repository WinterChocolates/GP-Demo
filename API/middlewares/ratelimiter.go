package middlewares

import (
	"net/http"
	"time"

	"API/storage/cache"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

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
