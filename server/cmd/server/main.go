package main

import (
	"fmt"
	"gin-vue/server/internal/config"
	"gin-vue/server/internal/router"
	"gin-vue/server/pkg/database"
)

func main() {
	// 1. 初始化配置
	config.Init()
	fmt.Println("✅ 配置加载成功")

	// 2. 初始化数据库
	if err := database.InitMySQL(); err != nil {
		panic(err)
	}
	fmt.Println("✅ 数据库连接成功")

	// 如果程序退出，延迟关闭数据库连接
	sqlDB, err := database.DB.DB()
	if err != nil {
		panic("获取底层sql.DB失败")
	}
	defer sqlDB.Close()

	// 3. 初始化路由
	r := router.Init()
	fmt.Println("✅ 路由初始化成功")

	// 4. 启动服务器
	serverPort := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	fmt.Printf("🚀 服务器即将启动于 http://127.0.0.1%s\n", serverPort)

	if err := r.Run(serverPort); err != nil {
		panic(fmt.Sprintf("服务器启动失败: %s", err))
	}
}
