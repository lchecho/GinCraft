package context

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	appCtx := New(nil)
	defer appCtx.Cancel()

	assert.NotEmpty(t, appCtx.GetTraceID())
	assert.False(t, appCtx.StartTime.IsZero())
	assert.NotNil(t, appCtx.CustomFields)
}

func TestSetUser(t *testing.T) {
	appCtx := New(nil)
	defer appCtx.Cancel()

	appCtx.SetUser("12345", "张三", "admin")

	assert.Equal(t, "12345", appCtx.GetUserID())
	assert.Equal(t, "张三", appCtx.GetUsername())
	assert.Equal(t, "admin", appCtx.GetUserRole())
}

func TestSetRequestInfo(t *testing.T) {
	appCtx := New(nil)
	defer appCtx.Cancel()

	appCtx.SetRequestInfo("GET", "/api/v1/users", "192.168.1.100", "Mozilla/5.0...")

	assert.Equal(t, "GET", appCtx.GetMethod())
	assert.Equal(t, "/api/v1/users", appCtx.GetPath())
	assert.Equal(t, "192.168.1.100", appCtx.GetClientIP())
	assert.Equal(t, "Mozilla/5.0...", appCtx.GetUserAgent())
}

func TestSetCustomField(t *testing.T) {
	appCtx := New(nil)
	defer appCtx.Cancel()

	appCtx.SetCustomField("test_key", "test_value")

	value, exists := appCtx.GetCustomField("test_key")
	assert.True(t, exists)
	assert.Equal(t, "test_value", value)

	// 测试不存在的字段
	_, exists = appCtx.GetCustomField("non_existent")
	assert.False(t, exists)
}

func TestGetCustomFields(t *testing.T) {
	appCtx := New(nil)
	defer appCtx.Cancel()

	appCtx.SetCustomField("key1", "value1")
	appCtx.SetCustomField("key2", "value2")

	fields := appCtx.GetCustomFields()
	assert.Len(t, fields, 2)
	assert.Equal(t, "value1", fields["key1"])
	assert.Equal(t, "value2", fields["key2"])
}

func TestGetDuration(t *testing.T) {
	appCtx := New(nil)
	defer appCtx.Cancel()

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	duration := appCtx.GetDuration()
	assert.Greater(t, duration, time.Duration(0))
}

func TestNewWithTimeout(t *testing.T) {
	appCtx := NewWithTimeout(nil, 100*time.Millisecond)
	defer appCtx.Cancel()

	// 等待超时
	time.Sleep(150 * time.Millisecond)

	assert.True(t, appCtx.IsCancelled())
}

func TestNewWithDeadline(t *testing.T) {
	deadline := time.Now().Add(100 * time.Millisecond)
	appCtx := NewWithDeadline(nil, deadline)
	defer appCtx.Cancel()

	// 等待截止时间
	time.Sleep(150 * time.Millisecond)

	assert.True(t, appCtx.IsCancelled())
}

func TestClone(t *testing.T) {
	originalCtx := New(nil)
	defer originalCtx.Cancel()

	originalCtx.SetUser("12345", "张三", "admin")
	originalCtx.SetCustomField("test_key", "test_value")

	clonedCtx := originalCtx.Clone()
	defer clonedCtx.Cancel()

	// 追踪ID应该相同
	assert.Equal(t, originalCtx.GetTraceID(), clonedCtx.GetTraceID())

	// 用户信息应该相同
	assert.Equal(t, originalCtx.GetUserID(), clonedCtx.GetUserID())
	assert.Equal(t, originalCtx.GetUsername(), clonedCtx.GetUsername())
	assert.Equal(t, originalCtx.GetUserRole(), clonedCtx.GetUserRole())

	// 自定义字段应该相同
	value, exists := clonedCtx.GetCustomField("test_key")
	assert.True(t, exists)
	assert.Equal(t, "test_value", value)
}

func TestContextInterface(t *testing.T) {
	appCtx := New(nil)
	defer appCtx.Cancel()

	// 测试context.Context接口方法
	_, ok := appCtx.Deadline()
	assert.False(t, ok) // 基础Context没有截止时间

	select {
	case <-appCtx.Done():
		t.Fatal("Context不应该被取消")
	default:
		// 正常情况
	}

	assert.Nil(t, appCtx.Err())
	assert.Nil(t, appCtx.Value("test_key"))
}

func TestLogging(t *testing.T) {
	// 创建一个测试logger
	logger := zap.NewNop()

	appCtx := New(nil)
	defer appCtx.Cancel()

	appCtx.SetLogger(logger)
	appCtx.SetUser("12345", "张三", "admin")
	appCtx.SetCustomField("test_key", "test_value")

	// 这些调用不应该panic
	appCtx.LogInfo("测试信息")
	appCtx.LogDebug("测试调试")
	appCtx.LogWarn("测试警告")
	appCtx.LogError("测试错误")
}

func TestConcurrency(t *testing.T) {
	customCtx := New(nil)
	defer customCtx.Cancel()

	// 启动多个goroutine同时访问Context
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			customCtx.SetCustomField("worker_id", id)
			customCtx.GetCustomFields()
			customCtx.GetDuration()
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证最终状态
	fields := customCtx.GetCustomFields()
	assert.Len(t, fields, 1) // 只有一个worker_id字段
}

func TestCancel(t *testing.T) {
	customCtx := New(nil)

	// 初始状态不应该被取消
	assert.False(t, customCtx.IsCancelled())

	// 取消Context
	customCtx.Cancel()

	// 现在应该被取消
	assert.True(t, customCtx.IsCancelled())
}
