// server/internal/router/middleware/log.go
package middleware

import (
	"bytes"
	"io"

	"server/internal/service"
	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// LogOperation 是一个操作日志中间件。
// 它会在每个API请求处理完毕后，异步地记录操作日志。
// 这样做可以避免日志记录阻塞主请求的响应。
func LogOperation(logService service.ILogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取 body 内容，因为 Gin 的 c.Request.Body 是一个只读流，读完即空。
		// 我们需要先读出来，再重新写回去，以便后续的 c.ShouldBindJSON 能正常工作。
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 先执行请求的业务逻辑
		c.Next()

		// 对用户操作进行日志记录
		// 为避免日志数据库充满大量无意义的查询记录，通常只记录状态变更的操作（非GET请求）
		if c.Request.Method == "GET" {
			return
		}

		// 从上下文中获取用户信息
		claims, exists := c.Get(ContextUserClaimsKey)
		if !exists {
			return // 获取不到用户信息，则不记录日志
		}

		userClaims, ok := claims.(*jwt.CustomClaims)
		if !ok {
			return // 用户信息格式不正确，不记录日志
		}

		// 调用服务，以goroutine方式异步创建日志记录
		go logService.CreateLog(
			userClaims,
			c.FullPath(),      // 路由路径作为模块名
			c.Request.Method,  // 请求方法作为操作类型
			string(bodyBytes), // 请求体作为参数
		)
	}
}
