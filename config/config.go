package config

import (
	"log"

	"github.com/spf13/viper"
)

// Init 初始化配置
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			log.Fatal("Error reading config file:", err)
		}
	}
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器配置
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	// 数据库配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "123456")
	viper.SetDefault("database.dbname", "matuto_blog")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parseTime", true)
	viper.SetDefault("database.loc", "Local")

	// 数据库连接池配置
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime_hours", 1)

	// JWT配置
	viper.SetDefault("jwt.secret", "matuto-blog-secret-key-change-in-production")
	viper.SetDefault("jwt.issuer", "matuto-blog")
	viper.SetDefault("jwt.access_token_ttl", 24)   // 小时
	viper.SetDefault("jwt.refresh_token_ttl", 168) // 小时 (7天)

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// 存储配置
	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local.base_path", "./uploads")
	viper.SetDefault("storage.local.base_url", "http://localhost:8080/uploads/")

	// 模板配置
	viper.SetDefault("theme.current", "default")
	viper.SetDefault("theme.path", "./web/templates")

	// CORS配置
	viper.SetDefault("cors.allowed_origins", "http://localhost:3000,http://localhost:8080")
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", "12h")
}

// GetString 获取字符串配置
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool 获取布尔配置
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetFloat64 获取浮点数配置
func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}
