package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// RespondSuccess 返回成功响应
func RespondWithJSON(c *gin.Context, statusCode int, code int, message string, data interface{}) {
	requestID, _ := c.Get("requestID")
	rid := ""
	if requestID != nil {
		rid = requestID.(string)
	}

	c.JSON(statusCode, StandardResponse{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: rid,
	})
}

func RespondSuccess(c *gin.Context, data interface{}) {
	RespondWithJSON(c, http.StatusOK, http.StatusOK, "成功", data)
}

func RespondCreated(c *gin.Context, data interface{}) {
	RespondWithJSON(c, http.StatusCreated, http.StatusCreated, "创建成功", data)
}

// RespondError 返回错误响应
func RespondError(c *gin.Context, status int, message string) {
	code := status
	if status == 0 {
		status = http.StatusInternalServerError
		code = http.StatusInternalServerError
	}
	RespondWithJSON(c, status, code, message, nil)
}

func MapStatusToString(status int) string {
	switch status {
	case http.StatusOK:
		return "available"
	case http.StatusServiceUnavailable:
		return "degraded"
	default:
		return "unknown"
	}
}
