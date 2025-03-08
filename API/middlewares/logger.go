package middlewares

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ResponseWriter 是对gin.ResponseWriter的包装，用于捕获响应内容
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法以捕获响应内容
func (w ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogger 基础请求日志中间件
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

// EnhancedLogger 增强版请求日志中间件
func EnhancedLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		protocol := c.Request.Proto
		userAgent := c.Request.UserAgent()
		remoteIP := c.ClientIP()
		requestID, _ := c.Get("requestID")

		// 记录请求开始
		logger.Info("请求开始",
			zap.String("request_id", requestID.(string)),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", remoteIP),
			zap.String("user_agent", userAgent),
			zap.String("protocol", protocol),
		)

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		// 根据状态码确定日志级别
		var logFunc func(msg string, fields ...zapcore.Field)
		if status >= 500 {
			logFunc = logger.Error
		} else if status >= 400 {
			logFunc = logger.Warn
		} else {
			logFunc = logger.Info
		}

		// 记录请求完成
		logFunc("请求完成",
			zap.String("request_id", requestID.(string)),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Int("size", size),
			zap.Duration("latency", latency),
			zap.String("ip", remoteIP),
		)
	}
}

// RequestBodyLogger 请求体日志中间件
func RequestBodyLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 仅记录POST/PUT/PATCH请求的请求体
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ := io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)

			requestID, _ := c.Get("requestID")
			logger.Debug("请求体",
				zap.String("request_id", requestID.(string)),
				zap.String("body", string(body)),
			)
		}
		c.Next()
	}
}

// ErrorLogger 错误日志中间件
func ErrorLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			requestID, _ := c.Get("requestID")
			for _, e := range c.Errors {
				logger.Error("请求处理错误",
					zap.String("request_id", requestID.(string)),
					zap.String("error", e.Error()),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
			}
		}
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

// EnhancedAuditLogger 增强版审计日志中间件
func EnhancedAuditLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			userID = "anonymous"
		}

		roles, _ := c.Get("roles")
		requestID, _ := c.Get("requestID")

		// 处理请求
		c.Next()

		// 记录审计日志
		logger.Info("审计记录",
			zap.String("request_id", requestID.(string)),
			zap.Any("user_id", userID),
			zap.Any("roles", roles),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
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
