package app

import (
	"fmt"

	"github.com/liuchen/gin-craft/internal/pkg/config"
	"github.com/liuchen/gin-craft/internal/pkg/cron"
	"github.com/liuchen/gin-craft/internal/pkg/database"
	"github.com/liuchen/gin-craft/pkg/logger"
	"go.uber.org/zap"
)

// Init 初始化应用
func Init(configPath string) error {
	// 加载配置
	if err := config.LoadConfig(configPath); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 初始化日志
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

	// 初始化数据库
	if err := database.InitDatabase(); err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// 初始化定时任务
	cron.InitCron()

	logger.Info("Application initialized successfully")
	return nil
}

// Close 关闭应用
func Close() {
	// 关闭数据库连接
	database.Close()

	// 停止定时任务
	cron.Stop()

	// 关闭日志
	logger.Close()
}
