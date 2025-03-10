package middlewares

import (
	"log"
	"net/http"
	"strings"

	"API/utils"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[JWT] 请求路径: %s %s | Authorization头: %s", c.Request.Method, c.Request.URL.Path, c.GetHeader("Authorization"))

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("[JWT] 错误: 请求头中未找到Authorization信息")
			utils.RespondError(c, http.StatusUnauthorized, "需要认证令牌")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Printf("[JWT] 解析失败: %v | 原始令牌: %s", err, tokenString[:10]+"...")
			utils.RespondError(c, http.StatusUnauthorized, "无效的令牌")
			c.Abort()
			return
		}

		log.Printf("[JWT] 解析成功 | 用户ID: %d | 角色: %v | 过期时间: %v", claims.UserID, claims.Roles, claims.ExpiresAt)

		c.Set("userID", uint(claims.UserID))
		c.Set("roles", claims.Roles)
		// 添加类型断言确保后续使用安全
		if claims.UserID == 0 {
			log.Println("[JWT] 错误: 无效的用户标识")
			utils.RespondError(c, http.StatusUnauthorized, "无效的用户标识")
			c.Abort()
			return
		}
		if c.Request.Method == "OPTIONS" {
			log.Println("[JWT] OPTIONS请求，跳过验证")
			c.Next()
			return
		}
		c.Next()
	}
}
