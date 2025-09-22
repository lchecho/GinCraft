package context

import (
	"time"

	"go.uber.org/zap"
)

// ExampleUsage 展示应用Context的使用示例
func ExampleUsage() {
	// 1. 创建基础Context
	appCtx := New(nil)
	defer appCtx.Cancel()

	// 2. 设置用户信息
	appCtx.SetUser("12345", "张三", "admin")

	// 3. 设置请求信息
	appCtx.SetRequestInfo("GET", "/api/v1/users", "192.168.1.100", "Mozilla/5.0...")

	// 4. 设置自定义字段
	appCtx.SetCustomField("business_type", "user_management")
	appCtx.SetCustomField("priority", "high")

	// 5. 记录日志
	appCtx.LogInfo("开始处理用户请求")
	appCtx.LogDebug("调试信息", zap.String("debug_key", "debug_value"))

	// 6. 模拟业务处理
	time.Sleep(100 * time.Millisecond)

	// 7. 记录处理结果
	appCtx.LogInfo("用户请求处理完成",
		zap.Duration("processing_time", appCtx.GetDuration()),
		zap.String("result", "success"),
	)

	// 8. 获取上下文信息
	traceID := appCtx.GetTraceID()
	userID := appCtx.GetUserID()
	duration := appCtx.GetDuration()

	// 记录获取到的信息
	appCtx.LogInfo("上下文信息",
		zap.String("trace_id", traceID),
		zap.String("user_id", userID),
		zap.Duration("duration", duration),
	)

	// 9. 获取自定义字段
	if businessType, exists := appCtx.GetCustomField("business_type"); exists {
		appCtx.LogInfo("业务类型", zap.String("business_type", businessType.(string)))
	}

	// 10. 克隆Context（用于子任务）
	subCtx := appCtx.Clone()
	defer subCtx.Cancel()

	// 子任务使用相同的追踪ID但不同的请求ID
	subCtx.LogInfo("子任务开始", zap.String("parent_trace_id", appCtx.GetTraceID()))
}

// ExampleWithTimeout 展示带超时的Context使用
func ExampleWithTimeout() {
	// 创建带5秒超时的Context
	appCtx := NewWithTimeout(nil, 5*time.Second)
	defer appCtx.Cancel()

	appCtx.SetUser("67890", "李四", "user")
	appCtx.SetCustomField("timeout_example", true)

	// 检查是否超时
	select {
	case <-appCtx.Done():
		appCtx.LogWarn("Context已超时或被取消")
	default:
		appCtx.LogInfo("Context仍在运行")
	}
}

// ExampleWithDeadline 展示带截止时间的Context使用
func ExampleWithDeadline() {
	// 创建带截止时间的Context（1分钟后）
	deadline := time.Now().Add(1 * time.Minute)
	appCtx := NewWithDeadline(nil, deadline)
	defer appCtx.Cancel()

	appCtx.SetUser("11111", "王五", "guest")
	appCtx.SetCustomField("deadline_example", true)

	// 获取截止时间
	if deadline, ok := appCtx.Deadline(); ok {
		appCtx.LogInfo("Context截止时间", zap.Time("deadline", deadline))
	}
}

// ExampleConcurrentUsage 展示并发使用Context
func ExampleConcurrentUsage() {
	// 创建主Context
	mainCtx := New(nil)
	defer mainCtx.Cancel()

	mainCtx.SetUser("99999", "赵六", "admin")
	mainCtx.SetCustomField("concurrent_example", true)

	// 启动多个goroutine，每个都使用克隆的Context
	for i := 0; i < 3; i++ {
		go func(workerID int) {
			workerCtx := mainCtx.Clone()
			defer workerCtx.Cancel()

			workerCtx.SetCustomField("worker_id", workerID)
			workerCtx.LogInfo("工作协程开始",
				zap.Int("worker_id", workerID),
				zap.String("parent_trace_id", mainCtx.GetTraceID()),
			)

			// 模拟工作
			time.Sleep(50 * time.Millisecond)

			workerCtx.LogInfo("工作协程完成",
				zap.Int("worker_id", workerID),
				zap.Duration("work_duration", workerCtx.GetDuration()),
			)
		}(i)
	}

	// 等待所有工作完成
	time.Sleep(100 * time.Millisecond)
	mainCtx.LogInfo("所有工作协程已完成")
}
