package database

import (
	"fmt"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// SQLiteDatabase SQLite数据库实现
type SQLiteDatabase struct {
	db       *gorm.DB
	filePath string
	mu       sync.RWMutex
}

// NewSQLiteDatabase 创建SQLite数据库实例
func NewSQLiteDatabase(filePath string) Database {
	if filePath == "" {
		filePath = "data.db" // 默认文件名
	}
	return &SQLiteDatabase{
		filePath: filePath,
	}
}

// Connect 连接SQLite数据库
func (s *SQLiteDatabase) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		return nil // 已经连接
	}

	var err error
	s.db, err = gorm.Open(sqlite.Open(s.filePath), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: NewGormLogger(), // 使用自定义日志
	})
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	return nil
}

// GetDB 获取数据库连接
func (s *SQLiteDatabase) GetDB() *gorm.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db
}

// Close 关闭数据库连接
func (s *SQLiteDatabase) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			return fmt.Errorf("failed to get DB instance: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close SQLite connection: %w", err)
		}
		s.db = nil
	}
	return nil
}

// Ping 测试数据库连接
func (s *SQLiteDatabase) Ping() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}

	return sqlDB.Ping()
}

// Migrate 数据库迁移
func (s *SQLiteDatabase) Migrate(models ...interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	return s.db.AutoMigrate(models...)
}
