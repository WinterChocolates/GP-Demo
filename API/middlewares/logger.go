package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// AuditLogger 审计日志中间件
func AuditLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logger.Info("audit",
			zap.String("user", c.GetString("userID")),
			zap.String("path", c.Request.URL.Path),
		)
	}
}

// OperationLog 操作日志中间件
func OperationLog(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logger.Info("operation",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("user", c.GetString("userID")),
		)
	}
}
