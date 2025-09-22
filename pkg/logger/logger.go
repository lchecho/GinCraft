package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log        *zap.Logger
	SugarLog   *zap.SugaredLogger
	TraceIDKey = "trace_id"
)

// InitLogger 初始化日志
func InitLogger(level, filename string, maxSize, maxBackups, maxAge int, compress bool) error {
	// 创建日志目录
	logDir := filepath.Dir(filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 设置日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置日志轮转
	hook := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize, // MB
		MaxBackups: maxBackups,
		MaxAge:     maxAge, // days
		Compress:   compress,
	}

	// 配置输出
	var core zapcore.Core
	// 同时输出到控制台和文件
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleOutput := zapcore.Lock(os.Stdout)
	consoleLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapLevel
	})

	// 文件输出
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	fileOutput := zapcore.AddSync(hook)
	fileLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapLevel
	})

	core = zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleOutput, consoleLevel),
		zapcore.NewCore(fileEncoder, fileOutput, fileLevel),
	)

	// 创建Logger
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	SugarLog = Log.Sugar()

	return nil
}

// WithTraceID 添加跟踪ID
func WithTraceID() string {
	return uuid.New().String()
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// Close 关闭日志
func Close() {
	_ = Log.Sync()
}

// Field 创建日志字段
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
