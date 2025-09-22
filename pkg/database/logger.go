package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/liuchen/gin-craft/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger 实现gorm的日志接口
type GormLogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

// NewGormLogger 创建gorm日志实例
func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold:         time.Second, // 慢查询阈值
		SkipErrRecordNotFound: true,        // 是否跳过记录未找到错误
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info 打印信息
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.Info(fmt.Sprintf(msg, data...))
}

// Warn 打印警告
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.Warn(fmt.Sprintf(msg, data...))
}

// Error 打印错误
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.Error(fmt.Sprintf(msg, data...))
}

// Trace 记录SQL执行
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 获取跟踪ID
	var traceID string
	if ctx != nil {
		if v, ok := ctx.Value(logger.TraceIDKey).(string); ok {
			traceID = v
		}
	}

	// 构建日志字段
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// 记录慢查询
	if elapsed > l.SlowThreshold {
		logger.Warn("SLOW SQL", fields...)
		return
	}

	// 记录错误
	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.SkipErrRecordNotFound) {
		fields = append(fields, zap.Error(err))
		logger.Error("SQL Error", fields...)
		return
	}

	// 记录正常查询
	logger.Debug("SQL", fields...)
}