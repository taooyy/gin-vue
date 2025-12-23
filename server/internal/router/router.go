package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	// 健康检查路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 后续的功能路由将在这里注册
	// apiGroup := r.Group("/api/v1")
	// {
	// 	 // 例如:
	//   // sys.RegisterUserRoutes(apiGroup.Group("/system"))
	// }

	return r
}
