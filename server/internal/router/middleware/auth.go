package middleware

import (
	"net/http"
	"strings"

	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const ContextUserClaimsKey = "userClaims"

// AuthMiddleware 创建一个 Gin 中间件，用于 JWT 认证
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "请求未携带 token，无权限访问",
			})
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token 格式不正确",
			})
			c.Abort()
			return
		}

		// parts[1] 是获取到的 tokenString，我们使用之前定义好的解析 JWT 的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的 Token",
			})
			c.Abort()
			return
		}

		// 将当前请求的 user信息 保存到请求的上下文 c 上
		c.Set(ContextUserClaimsKey, mc)
		c.Next() // 后续的处理函数可以用 c.Get("user") 来获取当前请求的用户信息
	}
}
