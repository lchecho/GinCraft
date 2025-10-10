package cron

import (
	"context"
	"fmt"
	"runtime/debug"

	customContext "github.com/liuchen/gin-craft/internal/pkg/context"
	"github.com/liuchen/gin-craft/pkg/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	Cron       *cron.Cron
	cronLogger *zap.Logger
)

// Job 定时任务接口
type Job interface {
	Execute(ctx *customContext.Context) error
	GetName() string
	GetDescription() string
}

// JobFunc 函数类型的定时任务
type JobFunc struct {
	name        string
	description string
	fn          func(ctx *customContext.Context) error
}

// NewJobFunc 创建函数类型的定时任务
func NewJobFunc(name, description string, fn func(ctx *customContext.Context) error) *JobFunc {
	return &JobFunc{
		name:        name,
		description: description,
		fn:          fn,
	}
}

func (j *JobFunc) Execute(ctx *customContext.Context) error {
	return j.fn(ctx)
}

func (j *JobFunc) GetName() string {
	return j.name
}

func (j *JobFunc) GetDescription() string {
	return j.description
}

// wrapJob 包装定时任务，添加trace_id、异常捕获和日志记录
func wrapJob(job Job) func() {
	return func() {
		// 创建自定义上下文
		ctx := customContext.New(context.Background())
		// 为上下文设置cron专用的logger
		ctx.SetLogger(logger.GetCronLogger())
		ctx.SetCustomField("job_name", job.GetName())
		ctx.SetCustomField("job_description", job.GetDescription())

		// 记录任务开始
		ctx.LogInfo("定时任务开始执行", zap.Time("start_time", ctx.StartTime))

		// 异常捕获和恢复
		defer func() {
			if r := recover(); r != nil {
				// 记录panic信息
				stack := debug.Stack()
				ctx.LogError("定时任务发生panic",
					zap.Any("panic", r),
					zap.String("stack", string(stack)),
					zap.Duration("duration", ctx.GetDuration()),
				)
			}
		}()

		// 执行任务
		err := job.Execute(ctx)

		// 记录任务结束
		duration := ctx.GetDuration()
		if err != nil {
			ctx.LogError("定时任务执行失败", zap.Error(err), zap.Duration("duration", duration))
		} else {
			ctx.LogInfo("定时任务执行成功", zap.Duration("duration", duration))
		}
	}
}

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
	cronLogger.Info("Cron调度器已启动")
}

// AddJob 添加定时任务
func AddJob(spec string, job Job) error {
	if Cron == nil {
		return fmt.Errorf("cron调度器未初始化")
	}

	wrappedJob := wrapJob(job)
	_, err := Cron.AddFunc(spec, wrappedJob)
	if err != nil {
		cronLogger.Error("添加定时任务失败",
			zap.String("job_name", job.GetName()),
			zap.String("spec", spec),
			zap.Error(err),
		)
		return fmt.Errorf("添加定时任务失败: %w", err)
	}

	cronLogger.Info("定时任务添加成功",
		zap.String("job_name", job.GetName()),
		zap.String("job_description", job.GetDescription()),
		zap.String("spec", spec),
	)

	return nil
}

// AddJobFunc 添加函数类型的定时任务
func AddJobFunc(spec, name, description string, fn func(ctx *customContext.Context) error) error {
	job := NewJobFunc(name, description, fn)
	return AddJob(spec, job)
}

// addJobs 添加定时任务
func addJobs() {
	// 这里可以添加你的定时任务
	// 示例：每分钟执行一次健康检查任务
	/*
		err := AddJobFunc("0 * * * * *", "health_check", "系统健康检查任务", func(ctx *customContext.Context) error {
			ctx.LogInfo("执行健康检查任务")
			// 在这里执行你的健康检查逻辑
			return nil
		})
		if err != nil {
			cronLogger.Error("添加健康检查任务失败", zap.Error(err))
		}
	*/

	// 示例：每小时执行一次数据清理任务
	/*
		err = AddJobFunc("0 0 * * * *", "data_cleanup", "数据清理任务", func(ctx *customContext.Context) error {
			ctx.LogInfo("执行数据清理任务")
			// 在这里执行你的数据清理逻辑
			return nil
		})
		if err != nil {
			cronLogger.Error("添加数据清理任务失败", zap.Error(err))
		}
	*/
}

// Stop 停止定时任务
func Stop() {
	if Cron != nil {
		Cron.Stop()
		if cronLogger != nil {
			cronLogger.Info("Cron调度器已停止")
		}
	}
}

// GetRunningJobs 获取正在运行的任务数量
func GetRunningJobs() int {
	if Cron == nil {
		return 0
	}
	return len(Cron.Entries())
}

// ListJobs 列出所有任务
func ListJobs() []cron.Entry {
	if Cron == nil {
		return nil
	}
	return Cron.Entries()
}
