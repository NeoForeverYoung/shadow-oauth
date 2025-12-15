package config

import (
	"os"
	"strconv"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string // 服务器端口
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string // SQLite 数据库文件路径
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string // JWT 签名密钥
	ExpireHours int    // Token 过期时间（小时）
}

// Load 加载配置，支持环境变量覆盖
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"), // 默认端口 8080
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./data/shadow.db"), // 默认数据库路径
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"), // 默认密钥（生产环境必须修改）
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),                          // 默认 24 小时
		},
	}

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

