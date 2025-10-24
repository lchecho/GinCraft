package context

import (
	"context"
	"github.com/liuchen/gin-craft/pkg/logger"
	"sync"
	"time"

	"go.uber.org/zap"
)

// CtxKey 用于在gin.Context中存储自定义Context的键
const CtxKey = "custom_context"

// Context 自定义上下文结构体，用于追踪请求或任务的上下文信息
type Context struct {
	// 原始context
	ctx context.Context
	// 取消函数
	cancel context.CancelFunc

	// 基础信息
	TraceID   string    // 追踪ID
	StartTime time.Time // 开始时间

	// 用户信息
	UserID   string // 用户ID
	Username string // 用户名
	UserRole string // 用户角色

	// 请求信息
	Method    string // HTTP方法
	Path      string // 请求路径
	ClientIP  string // 客户端IP
	UserAgent string // 用户代理

	// 自定义字段
	CustomFields map[string]interface{} // 自定义字段

	// 日志相关
	logger    *zap.Logger // 日志记录器
	logFields []zap.Field // 日志字段

	// 互斥锁，用于并发安全
	mu sync.RWMutex
}

// New 创建新的Context实例
func New(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.Background()
	}

	// 创建可取消的context
	ctx, cancel := context.WithCancel(ctx)

	return &Context{
		ctx:          ctx,
		cancel:       cancel,
		TraceID:      logger.WithTraceID(),
		StartTime:    time.Now(),
		CustomFields: make(map[string]interface{}),
		logFields:    make([]zap.Field, 0),
	}
}

// NewWithTimeout 创建带超时的Context实例
func NewWithTimeout(ctx context.Context, timeout time.Duration) *Context {
	if ctx == nil {
		ctx = context.Background()
	}

	// 创建带超时的context
	ctx, cancel := context.WithTimeout(ctx, timeout)

	return &Context{
		ctx:          ctx,
		cancel:       cancel,
		TraceID:      logger.WithTraceID(),
		StartTime:    time.Now(),
		CustomFields: make(map[string]interface{}),
		logFields:    make([]zap.Field, 0),
	}
}

// NewWithDeadline 创建带截止时间的Context实例
func NewWithDeadline(ctx context.Context, deadline time.Time) *Context {
	if ctx == nil {
		ctx = context.Background()
	}

	// 创建带截止时间的context
	ctx, cancel := context.WithDeadline(ctx, deadline)

	return &Context{
		ctx:          ctx,
		cancel:       cancel,
		TraceID:      logger.WithTraceID(),
		StartTime:    time.Now(),
		CustomFields: make(map[string]interface{}),
		logFields:    make([]zap.Field, 0),
	}
}

// GetContext 从gin.Context中获取应用Context
func GetContext(c context.Context) *Context {
	value := c.Value(CtxKey)
	if ctx, ok := value.(*Context); ok {
		return ctx
	}

	return nil
}

// MustGetContext 从头context.Context中获取应用Context，如果不存在则panic
func MustGetContext(c context.Context) *Context {
	ctx := GetContext(c)
	if ctx == nil {
		panic("App context not found")
	}

	return ctx
}

// SetUser 设置用户信息
func (c *Context) SetUser(userID, username, userRole string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.UserID = userID
	c.Username = username
	c.UserRole = userRole

	// 添加到日志字段
	c.AddLogField("user_id", userID)
	c.AddLogField("username", username)
	c.AddLogField("user_role", userRole)
}

// SetRequestInfo 设置请求信息
func (c *Context) SetRequestInfo(method, path, clientIP, userAgent string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Method = method
	c.Path = path
	c.ClientIP = clientIP
	c.UserAgent = userAgent

	// 添加到日志字段
	c.AddLogField("method", method)
	c.AddLogField("path", path)
	c.AddLogField("client_ip", clientIP)
	c.AddLogField("user_agent", userAgent)
}

// SetCustomField 设置自定义字段
func (c *Context) SetCustomField(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CustomFields[key] = value
	c.AddLogField(key, value)
}

// GetCustomField 获取自定义字段
func (c *Context) GetCustomField(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.CustomFields[key]
	return value, exists
}

// SetLogger 设置日志记录器
func (c *Context) SetLogger(logger *zap.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger = logger
}

// GetLogger 获取日志记录器
func (c *Context) GetLogger() *zap.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.logger
}

// LogInfo 记录信息日志
func (c *Context) LogInfo(msg string, fields ...zap.Field) {
	if c.logger != nil {
		allFields := c.buildLogFields(fields...)
		c.logger.Info(msg, allFields...)
	}
}

// LogWarn 记录警告日志
func (c *Context) LogWarn(msg string, fields ...zap.Field) {
	if c.logger != nil {
		allFields := c.buildLogFields(fields...)
		c.logger.Warn(msg, allFields...)
	}
}

// LogError 记录错误日志
func (c *Context) LogError(msg string, fields ...zap.Field) {
	if c.logger != nil {
		allFields := c.buildLogFields(fields...)
		c.logger.Error(msg, allFields...)
	}
}

// LogDebug 记录调试日志
func (c *Context) LogDebug(msg string, fields ...zap.Field) {
	if c.logger != nil {
		allFields := c.buildLogFields(fields...)
		c.logger.Debug(msg, allFields...)
	}
}

// GetDuration 获取上下文持续时间
func (c *Context) GetDuration() time.Duration {
	return time.Since(c.StartTime)
}

// GetTraceID 获取追踪ID
func (c *Context) GetTraceID() string {
	return c.TraceID
}

// GetUserID 获取用户ID
func (c *Context) GetUserID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.UserID
}

// GetUsername 获取用户名
func (c *Context) GetUsername() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Username
}

// GetUserRole 获取用户角色
func (c *Context) GetUserRole() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.UserRole
}

// GetMethod 获取HTTP方法
func (c *Context) GetMethod() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Method
}

// GetPath 获取请求路径
func (c *Context) GetPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Path
}

// GetClientIP 获取客户端IP
func (c *Context) GetClientIP() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ClientIP
}

// GetUserAgent 获取用户代理
func (c *Context) GetUserAgent() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.UserAgent
}

// GetCustomFields 获取所有自定义字段
func (c *Context) GetCustomFields() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 返回副本以避免并发问题
	result := make(map[string]interface{})
	for k, v := range c.CustomFields {
		result[k] = v
	}
	return result
}

// GetLogFields 获取所有日志字段
func (c *Context) GetLogFields() []zap.Field {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 返回副本以避免并发问题
	result := make([]zap.Field, len(c.logFields))
	copy(result, c.logFields)
	return result
}

// AddLogField 添加日志字段（内部方法，需要加锁）
func (c *Context) AddLogField(key string, value interface{}) {
	c.logFields = append(c.logFields, zap.Any(key, value))
}

// buildLogFields 构建日志字段（内部方法）
func (c *Context) buildLogFields(fields ...zap.Field) []zap.Field {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 基础字段
	baseFields := []zap.Field{
		zap.String("trace_id", c.TraceID),
		zap.Duration("duration", time.Since(c.StartTime)),
	}

	// 合并所有字段
	allFields := make([]zap.Field, 0, len(baseFields)+len(c.logFields)+len(fields))
	allFields = append(allFields, baseFields...)
	allFields = append(allFields, c.logFields...)
	allFields = append(allFields, fields...)

	return allFields
}

// Deadline 实现context.Context接口
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done 实现context.Context接口
func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err 实现context.Context接口
func (c *Context) Err() error {
	return c.ctx.Err()
}

// Value 实现context.Context接口
func (c *Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// Cancel 取消上下文
func (c *Context) Cancel() {
	c.cancel()
}

// IsCancelled 检查上下文是否已取消
func (c *Context) IsCancelled() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}

// Clone 克隆Context（创建新的实例但保留基础信息）
func (c *Context) Clone() *Context {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 创建新的context
	ctx, cancel := context.WithCancel(context.Background())

	clone := &Context{
		ctx:          ctx,
		cancel:       cancel,
		TraceID:      c.TraceID,  // 保持相同的追踪ID
		StartTime:    time.Now(), // 新的开始时间
		UserID:       c.UserID,
		Username:     c.Username,
		UserRole:     c.UserRole,
		Method:       c.Method,
		Path:         c.Path,
		ClientIP:     c.ClientIP,
		UserAgent:    c.UserAgent,
		CustomFields: make(map[string]interface{}),
		logger:       c.logger,
		logFields:    make([]zap.Field, 0),
	}

	// 复制自定义字段
	for k, v := range c.CustomFields {
		clone.CustomFields[k] = v
	}

	// 复制日志字段
	clone.logFields = make([]zap.Field, len(c.logFields))
	copy(clone.logFields, c.logFields)

	return clone
}
