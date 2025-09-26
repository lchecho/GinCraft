package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log *zap.Logger

	TraceIDKey = "trace_id"

	// 模块化日志管理器
	moduleLoggers = make(map[string]*ModuleLogger)
	loggerMutex   sync.RWMutex

	// 全局配置
	globalConfig LogConfig
)

// LogConfig 日志配置结构
type LogConfig struct {
	Level      string
	BaseDir    string // 日志基础目录
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// ModuleLogger 模块日志记录器
type ModuleLogger struct {
	name   string
	logger *zap.Logger
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
}

// InitLogger 初始化日志
func InitLogger(level, filename string, maxSize, maxBackups, maxAge int, compress bool) error {
	// 保存全局配置
	globalConfig = LogConfig{
		Level:      level,
		BaseDir:    filepath.Dir(filename),
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	// 创建日志目录
	if err := os.MkdirAll(globalConfig.BaseDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 初始化全局logger
	var err error
	Log, err = createLogger("app", filename, level, maxSize, maxBackups, maxAge, compress)
	if err != nil {
		return err
	}

	return nil
}

// createLogger 创建logger实例
func createLogger(name, filename, level string, maxSize, maxBackups, maxAge int, compress bool) (*zap.Logger, error) {
	// 设置日志级别
	zapLevel := parseLogLevel(level)

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

	// 创建Logger，为模块logger添加模块名字段
	options := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1)}
	if name != "app" {
		options = append(options, zap.Fields(zap.String("module", name)))
	}

	return zap.New(core, options...), nil
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
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

// GetModuleLogger 获取或创建模块日志记录器
func GetModuleLogger(moduleName string) Logger {
	loggerMutex.RLock()
	if logger, exists := moduleLoggers[moduleName]; exists {
		loggerMutex.RUnlock()
		return logger
	}
	loggerMutex.RUnlock()

	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	// 双重检查
	if logger, exists := moduleLoggers[moduleName]; exists {
		return logger
	}

	// 创建模块专用的日志文件名
	moduleFilename := filepath.Join(globalConfig.BaseDir, fmt.Sprintf("%s.log", moduleName))

	// 创建模块logger
	zapLogger, err := createLogger(
		moduleName,
		moduleFilename,
		globalConfig.Level,
		globalConfig.MaxSize,
		globalConfig.MaxBackups,
		globalConfig.MaxAge,
		globalConfig.Compress,
	)
	if err != nil {
		// 如果创建失败，返回默认logger
		return &moduleLoggerWrapper{Log}
	}

	moduleLogger := &ModuleLogger{
		name:   moduleName,
		logger: zapLogger,
	}

	moduleLoggers[moduleName] = moduleLogger
	return moduleLogger
}

// moduleLoggerWrapper 包装默认logger以实现Logger接口
type moduleLoggerWrapper struct {
	*zap.Logger
}

func (w *moduleLoggerWrapper) With(fields ...zap.Field) Logger {
	return &moduleLoggerWrapper{w.Logger.With(fields...)}
}

// ModuleLogger实现Logger接口
func (m *ModuleLogger) Debug(msg string, fields ...zap.Field) {
	m.logger.Debug(msg, fields...)
}

func (m *ModuleLogger) Info(msg string, fields ...zap.Field) {
	m.logger.Info(msg, fields...)
}

func (m *ModuleLogger) Warn(msg string, fields ...zap.Field) {
	m.logger.Warn(msg, fields...)
}

func (m *ModuleLogger) Error(msg string, fields ...zap.Field) {
	m.logger.Error(msg, fields...)
}

func (m *ModuleLogger) Fatal(msg string, fields ...zap.Field) {
	m.logger.Fatal(msg, fields...)
}

func (m *ModuleLogger) With(fields ...zap.Field) Logger {
	return &ModuleLogger{
		name:   m.name,
		logger: m.logger.With(fields...),
	}
}

// GetModuleName 获取模块名
func (m *ModuleLogger) GetModuleName() string {
	return m.name
}

// 预定义的模块logger获取函数
func GetDatabaseLogger() Logger {
	return GetModuleLogger("database")
}

func GetAPILogger() Logger {
	return GetModuleLogger("api")
}

func GetCacheLogger() Logger {
	return GetModuleLogger("cache")
}

func GetCronLogger() Logger {
	return GetModuleLogger("cron")
}

func GetNotificationLogger() Logger {
	return GetModuleLogger("notification")
}

// Close 关闭日志
func Close() {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	// 关闭所有模块logger
	for _, moduleLogger := range moduleLoggers {
		_ = moduleLogger.logger.Sync()
	}

	// 关闭全局logger
	if Log != nil {
		_ = Log.Sync()
	}
}

// Field 创建日志字段
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
