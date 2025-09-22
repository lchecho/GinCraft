package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

// MySQLDatabase MySQL数据库实现
type MySQLDatabase struct {
	db     *gorm.DB
	config *MySQLConfig
	mu     sync.RWMutex
}

// NewMySQLDatabase 创建MySQL数据库实例
func NewMySQLDatabase(config *MySQLConfig) Database {
	return &MySQLDatabase{
		config: config,
	}
}

// Connect 连接MySQL数据库
func (m *MySQLDatabase) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		return nil // 已经连接
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.config.Username,
		m.config.Password,
		m.config.Host,
		m.config.Port,
		m.config.Database,
	)

	var err error
	m.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: NewGormLogger(), // 使用自定义日志
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// 设置连接池
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}

	// 设置空闲连接池中的最大连接数
	sqlDB.SetMaxIdleConns(m.config.MaxIdleConns)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(m.config.MaxOpenConns)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Duration(m.config.ConnMaxLifetime) * time.Second)

	return nil
}

// GetDB 获取数据库连接
func (m *MySQLDatabase) GetDB() *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.db
}

// Close 关闭数据库连接
func (m *MySQLDatabase) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		sqlDB, err := m.db.DB()
		if err != nil {
			return fmt.Errorf("failed to get DB instance: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close MySQL connection: %w", err)
		}
		m.db = nil
	}
	return nil
}

// Ping 测试数据库连接
func (m *MySQLDatabase) Ping() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}

	return sqlDB.Ping()
}

// Migrate 数据库迁移
func (m *MySQLDatabase) Migrate(models ...interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	return m.db.AutoMigrate(models...)
}