package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// RespondSuccess 返回成功响应
func RespondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功",
		"data":    data,
	})
}

// RespondError 返回错误响应
func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"code":    1,
		"message": message,
		"data":    nil,
	})
}

func MapStatus(status int) string {
	switch status {
	case http.StatusOK:
		return "available"
	case http.StatusServiceUnavailable:
		return "degraded"
	default:
		return "unknown"
	}
}
