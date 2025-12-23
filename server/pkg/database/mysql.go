package database

import (
	"fmt"
	"gin-vue/server/internal/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL数据库连接
func InitMySQL() (err error) {
	// 从配置中获取 DSN
	dsn := config.Cfg.MySQL.Dsn()

	// GORM logger 配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound错误
			Colorful:                  true,        // 启用彩色打印
		},
	)

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return fmt.Errorf("无法连接到数据库: %w", err)
	}

	fmt.Println("数据库连接成功!")
	return nil
}
