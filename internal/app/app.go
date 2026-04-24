package app

import (
	"fmt"

	"github.com/liuchen/gin-craft/internal/pkg/config"
	"github.com/liuchen/gin-craft/internal/pkg/cron"
	"github.com/liuchen/gin-craft/internal/pkg/database"
	"github.com/liuchen/gin-craft/internal/pkg/redis"
	"github.com/liuchen/gin-craft/pkg/logger"
	"go.uber.org/zap"
)

// Init 初始化应用
func Init(configPath string) error {
	if err := config.LoadConfig(configPath); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := logger.InitLogger(
		config.Config.Log.Level,
		config.Config.Log.Filename,
		config.Config.Log.MaxSize,
		config.Config.Log.MaxBackups,
		config.Config.Log.MaxAge,
		config.Config.Log.Compress,
	); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := database.InitDatabase(); err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := redis.InitRedis(); err != nil {
		logger.Error("Failed to initialize Redis", zap.Error(err))
		// Redis 失败时回收已建立的 MySQL 连接，避免连接泄漏
		if closeErr := database.Close(); closeErr != nil {
			logger.Error("close database on rollback", zap.Error(closeErr))
		}
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}

	cron.InitCron()

	logger.Info("Application initialized successfully")
	return nil
}

// Close 关闭应用
func Close() {
	if err := database.Close(); err != nil {
		logger.Error("close database", zap.Error(err))
	}
	redis.Close()
	cron.Stop()
	logger.Close()
}
