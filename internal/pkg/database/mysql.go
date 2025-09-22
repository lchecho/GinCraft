package database

import (
	"sync"

	"github.com/liuchen/gin-craft/internal/pkg/config"
	"github.com/liuchen/gin-craft/pkg/database"
	"github.com/liuchen/gin-craft/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	db   database.Database
	once sync.Once
)

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	var err error
	once.Do(func() {
		cfg := config.Config.MySQL
		
		// 创建MySQL数据库实例
		mysqlConfig := &database.MySQLConfig{
			Host:            cfg.Host,
			Port:            cfg.Port,
			Username:        cfg.Username,
			Password:        cfg.Password,
			Database:        cfg.Database,
			MaxIdleConns:    cfg.MaxIdleConns,
			MaxOpenConns:    cfg.MaxOpenConns,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
		}
		
		db = database.NewMySQLDatabase(mysqlConfig)
		err = db.Connect()
		if err != nil {
			return
		}
		
		// 测试连接
		err = db.Ping()
		if err != nil {
			return
		}
		
		logger.Info("Database connected", zap.String("host", cfg.Host), zap.Int("port", cfg.Port))
	})
	return err
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	if db == nil {
		return nil
	}
	return db.GetDB()
}

// Close 关闭数据库连接
func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database", zap.Error(err))
		}
	}
}

// Migrate 数据库迁移
func Migrate(models ...interface{}) error {
	if db == nil {
		return nil
	}
	return db.Migrate(models...)
}
