package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB 全局数据库连接实例
	DB *gorm.DB
)

// Initialize 初始化数据库连接
func Initialize(dbPath string) error {
	// 确保数据库文件所在目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 连接到 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启 SQL 日志
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	DB = db
	log.Printf("数据库连接成功: %s", dbPath)
	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(models ...interface{}) error {
	if err := DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}
	log.Println("数据库表结构迁移成功")
	return nil
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
