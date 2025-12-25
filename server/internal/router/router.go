package router

import (
	"net/http"
	"server/internal/handler"
	"server/internal/repository"
	"server/internal/router/middleware"
	"server/internal/service"
	"server/pkg/database" // 确保导入 database 包

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
	userRepo := repository.NewUserRepository(database.DB)
	roleRepo := repository.NewRoleRepository(database.DB)
	orgRepo := repository.NewOrganizationRepository(database.DB)

	authService := service.NewAuthService(userRepo, roleRepo)
	accountService := service.NewAccountService(userRepo, roleRepo)
	schoolService := service.NewSchoolService(orgRepo, userRepo, roleRepo)

	authHandler := handler.NewAuthHandler(authService)
	accountHandler := handler.NewAccountHandler(accountService)
	schoolHandler := handler.NewSchoolHandler(schoolService)

	// --- 路由注册 ---
	apiGroup := r.Group("/api/v1")
	{
		// 登录路由，不需要认证
		sysGroup := apiGroup.Group("/system")
		{
			sysGroup.POST("/login", authHandler.Login)
		}

		// 账号管理路由，需要认证和授权
		accountGroup := apiGroup.Group("/accounts")
		accountGroup.Use(middleware.AuthMiddleware(), middleware.CanCreateUsers(roleRepo))
		{
			accountGroup.POST("", accountHandler.CreateAccount)
			accountGroup.GET("", accountHandler.ListAccounts)
			accountGroup.PUT("/:id/status", accountHandler.UpdateAccountStatus)
			accountGroup.DELETE("/:id", accountHandler.DeleteAccount)
			accountGroup.PUT("/:id", accountHandler.UpdateAccount)
			accountGroup.PUT("/:id/password", accountHandler.ResetPassword)
		}

		// 站点(学校)管理路由
		schoolGroup := apiGroup.Group("/schools")
		schoolGroup.Use(middleware.AuthMiddleware(), middleware.PlatformAdminAuth())
		{
			schoolGroup.POST("", schoolHandler.Create)
			schoolGroup.GET("", schoolHandler.List)
			schoolGroup.GET("/:id", schoolHandler.GetByID)
			schoolGroup.PUT("/:id", schoolHandler.Update)
			schoolGroup.DELETE("/:id", schoolHandler.Delete)
		}

		// 其他受保护的路由组
		protectedGroup := apiGroup.Group("")
		protectedGroup.Use(middleware.AuthMiddleware())
		{
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
