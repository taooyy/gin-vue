package router

import (
	"net/http"
	"server/internal/handler"
	"server/internal/repository"
	"server/internal/router/middleware"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

// Init 初始化并返回一个Gin引擎
func Init() *gin.Engine {
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// --- 依赖注入 ---
	// 在真实的应用中，这里可能会使用像 `wire` 或 `fx` 这样的依赖注入框架
	userRepo := repository.NewMockUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// --- 路由注册 ---
	apiGroup := r.Group("/api/v1")
	{
		// 系统相关路由
		sysGroup := apiGroup.Group("/system")
		{
			// 登录路由，不需要认证
			sysGroup.POST("/login", authHandler.Login)
		}

		// 受保护的路由组
		protectedGroup := apiGroup.Group("")
		protectedGroup.Use(middleware.AuthMiddleware())
		{
			// 在这里注册所有需要 JWT 认证的路由
			protectedGroup.GET("/ping-auth", func(c *gin.Context) {
				claims, _ := c.Get(middleware.ContextUserClaimsKey)
				c.JSON(http.StatusOK, gin.H{
					"message": "pong from authenticated user",
					"user":    claims,
				})
			})
		}
	}

	// 健康检查路由 (公开)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}
