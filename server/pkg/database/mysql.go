package database

import (
	"fmt"
	"log"
	"os"
	"server/internal/config"
	"server/internal/model" // 导入所有模型的包
	"server/pkg/password"
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
	fmt.Println("正在检查并执行数据迁移...")
	err := DB.AutoMigrate(
		// System models
		&model.SysOrganization{},
		&model.SysUser{},
		&model.SysRole{}, // <-- 添加 SysRole 模型
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
		return fmt.Errorf("数据库迁移失败: %w", err)
	}
	fmt.Println("✅ 数据模型迁移成功！")

	// 执行数据填充
	if err := seedRoles(DB); err != nil {
		return fmt.Errorf("角色数据填充失败: %w", err)
	}
	if err := seedUsers(DB); err != nil {
		return fmt.Errorf("初始用户填充失败: %w", err)
	}

	return nil
}

// seedRoles 填充初始的角色数据
func seedRoles(db *gorm.DB) error {
	roles := []model.SysRole{
		{RoleName: "平台管理员", RoleKey: "platform_admin", CanCreateUsers: true},
		{RoleName: "平台员工", RoleKey: "platform_staff", CanCreateUsers: false},
		{RoleName: "学校管理员", RoleKey: "school_admin", CanCreateUsers: true},
		{RoleName: "学校员工", RoleKey: "school_staff", CanCreateUsers: false},
		{RoleName: "供应商管理员", RoleKey: "supplier_admin", CanCreateUsers: true},
		{RoleName: "供应商员工", RoleKey: "supplier_staff", CanCreateUsers: false},
	}

	fmt.Println("正在填充初始角色数据...")
	for _, role := range roles {
		// 如果记录不存在则创建
		result := db.Where(model.SysRole{RoleKey: role.RoleKey}).FirstOrCreate(&role)
		if result.Error != nil {
			return result.Error
		}
	}

	fmt.Println("✅ 角色数据填充成功！")
	return nil
}

// seedUsers 填充初始的用户数据
func seedUsers(db *gorm.DB) error {
	fmt.Println("正在填充初始用户数据...")
	// 1. 获取 "平台管理员" 的角色ID
	var adminRole model.SysRole
	if err := db.Where("role_key = ?", "platform_admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("找不到 'platform_admin' 角色: %w", err)
	}

	// 2. 为初始用户密码进行哈希
	hashedPassword, err := password.Hash("password123")
	if err != nil {
		return fmt.Errorf("初始密码哈希失败: %w", err)
	}

	// 3. 定义初始管理员用户
	adminUser := model.SysUser{
		Username: "platform_admin",
		Password: hashedPassword,
		RealName: "平台初始管理员",
		RoleID:   adminRole.ID,
		Status:   1, // 正常状态
		// OrgID 可以根据需要设置，这里默认为 0 或 1
		OrgID: 1,
	}

	// 4. 如果用户不存在则创建
	result := db.Where(model.SysUser{Username: adminUser.Username}).FirstOrCreate(&adminUser)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		fmt.Println("✅ 初始用户 [platform_admin] 创建成功！")
	} else {
		fmt.Println("初始用户 [platform_admin] 已存在，跳过创建。")
	}

	return nil
}
