package middlewares

import (
	"net/http"

	"API/utils"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			utils.RespondError(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		hasAdmin := false
		for _, role := range roles.([]string) {
			if role == "admin" {
				hasAdmin = true
				break
			}
		}

		if !hasAdmin {
			utils.RespondError(c, http.StatusForbidden, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
