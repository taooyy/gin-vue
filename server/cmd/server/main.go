package main

import (
	"flag"
	"fmt"
	"server/internal/config"
	"server/internal/router"
	"server/pkg/database"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化环境和配置
	var env string
	flag.StringVar(&env, "env", "dev", "设置运行环境, 可选值为 dev, test, prod")
	flag.Parse()

	configName := fmt.Sprintf("config.%s", env)
	config.Init(configName)
	fmt.Printf("✅ [%s] 环境配置加载成功\n", env)

	// 根据环境设置 Gin 模式
	switch strings.ToLower(env) {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 2. 初始化数据库
	if err := database.InitMySQL(); err != nil {
		panic(err)
	}
	// 执行数据迁移
	if err := database.Migrate(); err != nil {
		panic(fmt.Sprintf("数据迁移失败: %s", err))
	}
	fmt.Printf("✅ [%s] 环境数据库连接成功 & 数据迁移完成\n", env)

	// 如果程序退出，延迟关闭数据库连接
	sqlDB, err := database.DB.DB()
	if err != nil {
		panic("获取底层sql.DB失败")
	}
	defer sqlDB.Close()

	// 3. 初始化路由
	r := router.Init()
	fmt.Printf("✅ [%s] 环境路由初始化成功\n", env)

	// 4. 启动服务器
	serverPort := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	fmt.Printf("🚀 服务器即将在 [%s] 环境启动, Gin 模式: [%s], 监听端口: [http://127.0.0.1%s]\n", env, gin.Mode(), serverPort)

	if err := r.Run(serverPort); err != nil {
		panic(fmt.Sprintf("服务器启动失败: %s", err))
	}
}
