// server/internal/router/middleware/authz.go
package middleware

import (
	"net/http"

	"server/internal/repository"
	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// CanCreateUsers 是一个授权中间件，用于检查当前用户是否有权限创建新用户
// 它依赖 IRoleRepository 来从数据库查询角色的权限
func CanCreateUsers(roleRepo repository.IRoleRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Gin 上下文中获取由 AuthMiddleware 放入的 claims
		claims, exists := c.Get(ContextUserClaimsKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "无法获取用户信息，禁止访问"})
			c.Abort()
			return
		}

		// 2. 断言 claims 的类型
		userClaims, ok := claims.(*jwt.CustomClaims)
		if !ok || userClaims == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "用户信息格式错误，禁止访问"})
			c.Abort()
			return
		}

		// 3. 从 claims 中获取角色标识 (roleKey)
		roleKey := userClaims.Role
		if roleKey == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "用户角色未定义，禁止访问"})
			c.Abort()
			return
		}

		// 4. 使用 roleRepo 查询数据库，获取角色的详细信息
		role, err := roleRepo.FindRoleByRoleKey(roleKey)
		if err != nil {
			// 如果在数据库中找不到角色，也视为无权限
			c.JSON(http.StatusForbidden, gin.H{"error": "角色不存在或查询失败，禁止访问"})
			c.Abort()
			return
		}

		// 5. 检查角色的 CanCreateUsers 权限
		if !role.CanCreateUsers {
			c.JSON(http.StatusForbidden, gin.H{"error": "您没有创建新账号的权限"})
			c.Abort()
			return
		}

		// 6. 权限检查通过，放行请求
		c.Next()
	}
}

const RolePlatformAdmin = "platform_admin"

const RolePlatformStaff = "platform_staff"

// PlatformAdminAuth 是一个授权中间件，用于检查当前用户是否为平台管理员或员工
func PlatformAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(ContextUserClaimsKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "无法获取用户信息，禁止访问"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*jwt.CustomClaims)
		if !ok || userClaims == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "用户信息格式错误，禁止访问"})
			c.Abort()
			return
		}

		// 检查角色是否为平台管理员或员工
		if userClaims.Role != RolePlatformAdmin && userClaims.Role != RolePlatformStaff {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要平台管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
