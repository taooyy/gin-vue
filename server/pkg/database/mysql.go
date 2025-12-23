package database

import (
	"fmt"
	"log"
	"os"
	"server/internal/config"
	"server/internal/model" // 导入所有模型的包
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
	fmt.Println("🚀 即将开始数据迁移...")
	return nil
}

// Migrate 执行数据迁移
func Migrate() error {
	// 检查核心表是否存在，如果存在则跳过迁移
	if DB.Migrator().HasTable(&model.SysUser{}) {
		fmt.Println("数据库表已存在，跳过迁移。")
		return nil
	}

	fmt.Println("正在进行首次数据迁移...")
	err := DB.AutoMigrate(
		// System models
		&model.SysOrganization{},
		&model.SysUser{},
		&model.SysDictionary{},
		&model.SysOpLog{},
		&model.SysBanner{},

		// SCM models
		&model.ScmCategory{},
		&model.ScmProduct{},
		&model.ScmProductQuote{},
		&model.ScmSupplierStaff{},

		// Order models
		&model.OrdCart{},
		&model.OrdOrder{},
		&model.OrdOrderItem{},
		&model.OrdAfterSale{},
		&model.OrdItemTrace{},

		// Finance models
		&model.FinBill{},
		&model.FinStatement{},
	)
	if err != nil {
		return err
	}
	fmt.Println("✅ 首次数据迁移成功！")
	return nil
}
