package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestContext 中间件为每个请求添加唯一ID和上下文管理
func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求唯一ID
		requestID := uuid.New().String()
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)

		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		// 将上下文附加到请求
		c.Request = c.Request.WithContext(ctx)

		// 继续处理请求
		c.Next()
	}
}

// PerformanceMonitor 中间件用于监控API性能
func PerformanceMonitor(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)

		// 记录慢请求
		if latency > 200*time.Millisecond {
			requestID, _ := c.Get("requestID")
			logger.Warn("慢请求检测",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Duration("latency", latency),
				zap.String("requestID", requestID.(string)),
			)
		}

		// 添加性能指标到响应头
		c.Header("X-Response-Time", latency.String())
	}
}
