package database

import (
	"sync"

	"github.com/liuchen/gin-craft/internal/pkg/config"
	pkgdb "github.com/liuchen/gin-craft/pkg/database"
	"gorm.io/gorm"
)

var (
	db   pkgdb.Database
	once sync.Once
)

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	var err error
	once.Do(func() {
		cfg := config.Config.MySQL
		mysqlConfig := &pkgdb.MySQLConfig{
			Host:            cfg.Host,
			Port:            cfg.Port,
			Username:        cfg.Username,
			Password:        cfg.Password,
			Database:        cfg.Database,
			MaxIdleConns:    cfg.MaxIdleConns,
			MaxOpenConns:    cfg.MaxOpenConns,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
		}

		db = pkgdb.NewMySQLDatabase(mysqlConfig)
		if err = db.Connect(); err != nil {
			return
		}
		err = db.Ping()
	})
	return err
}

// GetDatabase 获取 Database 接口实例
func GetDatabase() pkgdb.Database {
	return db
}

// GetDB 获取底层 *gorm.DB
func GetDB() *gorm.DB {
	if db == nil {
		return nil
	}
	return db.GetDB()
}

// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

// Migrate 数据库迁移
func Migrate(models ...interface{}) error {
	if db == nil {
		return nil
	}
	return db.Migrate(models...)
}
