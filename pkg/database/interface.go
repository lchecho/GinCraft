package database

import (
	"gorm.io/gorm"
)

// Database 数据库接口
type Database interface {
	// Connect 连接数据库
	Connect() error
	// GetDB 获取数据库连接
	GetDB() *gorm.DB
	// Close 关闭数据库连接
	Close() error
	// Ping 测试数据库连接
	Ping() error
	// Migrate 数据库迁移
	Migrate(models ...interface{}) error
}