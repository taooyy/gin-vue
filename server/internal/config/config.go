package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 是应用的总配置结构体
var Cfg AppConfig

type AppConfig struct {
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Server ServerConfig `mapstructure:"server"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	Config   string `mapstructure:"config"`
}

// Dsn 返回格式化的数据库连接字符串
func (m *MySQLConfig) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		m.User, m.Password, m.Host, m.Port, m.DbName, m.Config)
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// Init 初始化配置
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs") // 定义配置文件的路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	// 将配置解析到 Cfg 变量
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(fmt.Errorf("解析配置到结构体失败: %w", err))
	}
}
