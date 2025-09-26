package cron

import (
	"github.com/liuchen/gin-craft/pkg/logger"
	"github.com/robfig/cron/v3"
)

var (
	Cron       *cron.Cron
	cronLogger logger.Logger
)

// InitCron 初始化定时任务
func InitCron() {
	// 初始化cron模块logger
	cronLogger = logger.GetCronLogger()

	// 创建定时任务调度器，支持秒级别的定时任务
	Cron = cron.New(cron.WithSeconds())

	// 添加定时任务
	addJobs()

	// 启动定时任务
	Cron.Start()
	cronLogger.Info("Cron scheduler started")
}

// addJobs 添加定时任务
func addJobs() {
	// 这里可以添加你的定时任务
	// 示例：每分钟执行一次
	// _, err := Cron.AddFunc("0 * * * * *", func() {
	//     cronLogger.Info("Running cron job")
	//     // 在这里执行你的定时任务
	// })
	// if err != nil {
	//     cronLogger.Error("Failed to add cron job", zap.Error(err))
	// }

	// 添加更多定时任务...
}

// Stop 停止定时任务
func Stop() {
	if Cron != nil {
		Cron.Stop()
		if cronLogger != nil {
			cronLogger.Info("Cron scheduler stopped")
		}
	}
}
