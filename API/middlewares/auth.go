package middlewares

import (
	"net/http"
	"strings"

	"API/utils"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.RespondError(c, http.StatusUnauthorized, "需要认证令牌")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "无效的令牌")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}
