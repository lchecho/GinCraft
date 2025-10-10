package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// 初始化logger
	err := InitLogger("debug", logFile, 10, 3, 7, false)
	require.NoError(t, err)

	// 验证全局logger被创建
	assert.NotNil(t, Log)

	// 验证配置被保存
	assert.Equal(t, "debug", globalConfig.Level)
	assert.Equal(t, tmpDir, globalConfig.BaseDir)
	assert.Equal(t, 10, globalConfig.MaxSize)
	assert.Equal(t, 3, globalConfig.MaxBackups)
	assert.Equal(t, 7, globalConfig.MaxAge)
	assert.False(t, globalConfig.Compress)

	// 清理
	Close()
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"debug", "debug"},
		{"DEBUG", "debug"},
		{"info", "info"},
		{"INFO", "info"},
		{"warn", "warn"},
		{"WARN", "warn"},
		{"error", "error"},
		{"ERROR", "error"},
		{"fatal", "fatal"},
		{"FATAL", "fatal"},
		{"panic", "panic"},
		{"PANIC", "panic"},
		{"unknown", "info"}, // 默认级别
		{"", "info"},        // 空字符串默认级别
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			level := parseLogLevel(tt.input)
			assert.Equal(t, tt.expected, level.String())
		})
	}
}

func TestGetModuleLogger(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// 初始化logger
	err := InitLogger("info", logFile, 10, 3, 7, false)
	require.NoError(t, err)

	// 获取模块logger
	dbLogger := GetModuleLogger("database")
	assert.NotNil(t, dbLogger)

	// 再次获取同一个模块logger，应该返回同一个实例
	dbLogger2 := GetModuleLogger("database")
	assert.Equal(t, dbLogger, dbLogger2)

	// 获取不同的模块logger
	authLogger := GetModuleLogger("auth")
	assert.NotNil(t, authLogger)
	assert.NotEqual(t, dbLogger, authLogger)

	// 验证模块日志文件被创建
	dbLogFile := filepath.Join(tmpDir, "database.log")
	authLogFile := filepath.Join(tmpDir, "auth.log")

	// 写入一些日志
	dbLogger.Info("Database connection established")
	authLogger.Info("User authenticated")

	// 等待日志写入
	time.Sleep(100 * time.Millisecond)

	// 验证文件存在
	_, err = os.Stat(dbLogFile)
	assert.NoError(t, err)
	_, err = os.Stat(authLogFile)
	assert.NoError(t, err)

	// 清理
	Close()
}

func TestPredefinedModuleLoggers(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// 初始化logger
	err := InitLogger("debug", logFile, 10, 3, 7, false)
	require.NoError(t, err)

	// 测试预定义的模块logger
	tests := []struct {
		name   string
		getter func() *zap.Logger
	}{
		{"database", GetDatabaseLogger},
		{"api", GetAPILogger},
		{"cache", GetCacheLogger},
		{"cron", GetCronLogger},
		{"notification", GetNotificationLogger},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := tt.getter()
			assert.NotNil(t, logger)

			// 测试日志方法
			logger.Debug("Debug message")
			logger.Info("Info message")
			logger.Warn("Warn message")
			logger.Error("Error message")

			// 测试With方法
			contextLogger := logger.With(zap.String("test_field", "test_value"))
			assert.NotNil(t, contextLogger)
			contextLogger.Info("Context message")
		})
	}

	// 清理
	Close()
}

func TestModuleLoggerInterface(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// 初始化logger
	err := InitLogger("debug", logFile, 10, 3, 7, false)
	require.NoError(t, err)

	// 获取模块logger
	logger := GetModuleLogger("test")

	// 测试所有接口方法
	logger.Debug("Debug message", zap.String("key", "value"))
	logger.Info("Info message", zap.String("key", "value"))
	logger.Warn("Warn message", zap.String("key", "value"))
	logger.Error("Error message", zap.String("key", "value"))

	// 测试With方法返回的也是Logger接口
	contextLogger := logger.With(zap.String("context", "test"))
	contextLogger.Info("Context message")

	// 清理
	Close()
}

func TestModuleLoggerWithFields(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// 初始化logger
	err := InitLogger("debug", logFile, 10, 3, 7, false)
	require.NoError(t, err)

	// 获取模块logger
	logger := GetModuleLogger("test")

	// 测试带字段的日志记录
	logger.Info("Test message",
		zap.String("string_field", "value"),
		zap.Int("int_field", 123),
		zap.Bool("bool_field", true),
	)

	// 测试With方法链式调用
	contextLogger := logger.With(
		zap.String("request_id", "req-123"),
		zap.String("user_id", "user-456"),
	)

	contextLogger.Info("Context message",
		zap.String("action", "login"),
		zap.Duration("duration", 100*time.Millisecond),
	)

	// 清理
	Close()
}

func TestConcurrentModuleLoggerAccess(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// 初始化logger
	err := InitLogger("info", logFile, 10, 3, 7, false)
	require.NoError(t, err)

	// 并发获取同一个模块logger
	done := make(chan *zap.Logger, 10)
	moduleName := "concurrent_test"

	for i := 0; i < 10; i++ {
		go func() {
			logger := GetModuleLogger(moduleName)
			logger.Info("Concurrent message")
			done <- logger
		}()
	}

	// 收集所有logger实例
	loggers := make([]*zap.Logger, 10)
	for i := 0; i < 10; i++ {
		loggers[i] = <-done
	}

	// 验证所有实例都是同一个
	for i := 1; i < len(loggers); i++ {
		assert.Equal(t, loggers[0], loggers[i])
	}

	// 清理
	Close()
}

func TestField(t *testing.T) {
	field := Field("test_key", "test_value")
	assert.Equal(t, "test_key", field.Key)
	assert.Equal(t, "test_value", field.String)
}

func TestClose(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// 初始化logger
	err := InitLogger("debug", logFile, 10, 3, 7, false)
	require.NoError(t, err)

	// 记录初始的模块logger数量
	initialCount := len(moduleLoggers)

	// 创建一些模块logger
	GetModuleLogger("close_test1")
	GetModuleLogger("close_test2")
	GetModuleLogger("close_test3")

	// 验证模块logger被创建
	assert.Len(t, moduleLoggers, initialCount+3)

	// 关闭所有logger
	Close()

	// 注意：Close()不会清空moduleLoggers映射，只是同步日志
	// 这是正确的行为，因为logger实例仍然可以使用
	// 验证映射仍然包含我们创建的logger
	assert.Len(t, moduleLoggers, initialCount+3)
}
