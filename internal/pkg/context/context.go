package context

import (
	stdctx "context"
	"sync"
	"time"

	"github.com/liuchen/gin-craft/pkg/logger"
	"go.uber.org/zap"
)

// ctxKeyType 未导出的 context key 类型，避免 key 碰撞
type ctxKeyType struct{}

// CtxKey 自定义 context key（在 context.WithValue 中使用）
var CtxKey = ctxKeyType{}

// Context 自定义上下文结构体
type Context struct {
	ctx    stdctx.Context
	cancel stdctx.CancelFunc

	TraceID   string
	StartTime time.Time

	UserID   string
	Username string
	UserRole string

	Method    string
	Path      string
	ClientIP  string
	UserAgent string

	CustomFields map[string]interface{}

	logger    *zap.Logger
	logFields []zap.Field

	mu sync.RWMutex
}

// New 创建新的 Context 实例
func New(ctx stdctx.Context) *Context {
	if ctx == nil {
		ctx = stdctx.Background()
	}
	ctx, cancel := stdctx.WithCancel(ctx)
	return &Context{
		ctx:          ctx,
		cancel:       cancel,
		TraceID:      logger.WithTraceID(),
		StartTime:    time.Now(),
		CustomFields: make(map[string]interface{}),
		logFields:    make([]zap.Field, 0),
	}
}

// NewWithTraceID 使用指定 traceID 创建 Context
func NewWithTraceID(ctx stdctx.Context, traceID string) *Context {
	c := New(ctx)
	if traceID != "" {
		c.TraceID = traceID
	}
	return c
}

// NewWithTimeout 创建带超时的 Context
func NewWithTimeout(ctx stdctx.Context, timeout time.Duration) *Context {
	if ctx == nil {
		ctx = stdctx.Background()
	}
	ctx, cancel := stdctx.WithTimeout(ctx, timeout)
	return &Context{
		ctx:          ctx,
		cancel:       cancel,
		TraceID:      logger.WithTraceID(),
		StartTime:    time.Now(),
		CustomFields: make(map[string]interface{}),
		logFields:    make([]zap.Field, 0),
	}
}

// NewWithDeadline 创建带截止时间的 Context
func NewWithDeadline(ctx stdctx.Context, deadline time.Time) *Context {
	if ctx == nil {
		ctx = stdctx.Background()
	}
	ctx, cancel := stdctx.WithDeadline(ctx, deadline)
	return &Context{
		ctx:          ctx,
		cancel:       cancel,
		TraceID:      logger.WithTraceID(),
		StartTime:    time.Now(),
		CustomFields: make(map[string]interface{}),
		logFields:    make([]zap.Field, 0),
	}
}

// GetContext 从任意 context.Context 中取出自定义 Context
func GetContext(c stdctx.Context) *Context {
	if c == nil {
		return nil
	}
	if v, ok := c.Value(CtxKey).(*Context); ok {
		return v
	}
	return nil
}

// MustGetContext 同 GetContext，找不到则 panic
func MustGetContext(c stdctx.Context) *Context {
	ctx := GetContext(c)
	if ctx == nil {
		panic("App context not found")
	}
	return ctx
}

// SetUser 设置用户信息
func (c *Context) SetUser(userID, username, userRole string) {
	c.mu.Lock()
	c.UserID = userID
	c.Username = username
	c.UserRole = userRole
	c.logFields = append(c.logFields,
		zap.String("user_id", userID),
		zap.String("username", username),
		zap.String("user_role", userRole),
	)
	c.mu.Unlock()
}

// SetRequestInfo 设置请求信息
func (c *Context) SetRequestInfo(method, path, clientIP, userAgent string) {
	c.mu.Lock()
	c.Method = method
	c.Path = path
	c.ClientIP = clientIP
	c.UserAgent = userAgent
	c.logFields = append(c.logFields,
		zap.String("method", method),
		zap.String("path", path),
		zap.String("client_ip", clientIP),
		zap.String("user_agent", userAgent),
	)
	c.mu.Unlock()
}

// SetCustomField 设置自定义字段
func (c *Context) SetCustomField(key string, value interface{}) {
	c.mu.Lock()
	c.CustomFields[key] = value
	c.logFields = append(c.logFields, zap.Any(key, value))
	c.mu.Unlock()
}

// GetCustomField 获取自定义字段
func (c *Context) GetCustomField(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.CustomFields[key]
	return v, ok
}

// SetLogger 设置日志记录器
func (c *Context) SetLogger(l *zap.Logger) {
	c.mu.Lock()
	c.logger = l
	c.mu.Unlock()
}

// GetLogger 获取日志记录器
func (c *Context) GetLogger() *zap.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}

// LogInfo 记录信息日志
func (c *Context) LogInfo(msg string, fields ...zap.Field) {
	if l := c.GetLogger(); l != nil {
		l.Info(msg, c.buildLogFields(fields...)...)
	}
}

// LogWarn 记录警告日志
func (c *Context) LogWarn(msg string, fields ...zap.Field) {
	if l := c.GetLogger(); l != nil {
		l.Warn(msg, c.buildLogFields(fields...)...)
	}
}

// LogError 记录错误日志
func (c *Context) LogError(msg string, fields ...zap.Field) {
	if l := c.GetLogger(); l != nil {
		l.Error(msg, c.buildLogFields(fields...)...)
	}
}

// LogDebug 记录调试日志
func (c *Context) LogDebug(msg string, fields ...zap.Field) {
	if l := c.GetLogger(); l != nil {
		l.Debug(msg, c.buildLogFields(fields...)...)
	}
}

// GetDuration 获取上下文持续时间
func (c *Context) GetDuration() time.Duration {
	return time.Since(c.StartTime)
}

// GetTraceID 获取追踪 ID
func (c *Context) GetTraceID() string {
	return c.TraceID
}

// GetUserID 获取用户 ID
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

// GetMethod 获取 HTTP 方法
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

// GetClientIP 获取客户端 IP
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

// GetCustomFields 获取所有自定义字段（副本）
func (c *Context) GetCustomFields() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]interface{}, len(c.CustomFields))
	for k, v := range c.CustomFields {
		result[k] = v
	}
	return result
}

// GetLogFields 获取所有日志字段（副本）
func (c *Context) GetLogFields() []zap.Field {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]zap.Field, len(c.logFields))
	copy(result, c.logFields)
	return result
}

// AddLogField 并发安全地追加日志字段
func (c *Context) AddLogField(key string, value interface{}) {
	c.mu.Lock()
	c.logFields = append(c.logFields, zap.Any(key, value))
	c.mu.Unlock()
}

func (c *Context) buildLogFields(extra ...zap.Field) []zap.Field {
	c.mu.RLock()
	defer c.mu.RUnlock()

	base := []zap.Field{
		zap.String("trace_id", c.TraceID),
		zap.Duration("duration", time.Since(c.StartTime)),
	}
	all := make([]zap.Field, 0, len(base)+len(c.logFields)+len(extra))
	all = append(all, base...)
	all = append(all, c.logFields...)
	all = append(all, extra...)
	return all
}

// ----- context.Context 接口实现 -----

func (c *Context) Deadline() (time.Time, bool)    { return c.ctx.Deadline() }
func (c *Context) Done() <-chan struct{}          { return c.ctx.Done() }
func (c *Context) Err() error                     { return c.ctx.Err() }
func (c *Context) Value(key interface{}) interface{} {
	if key == CtxKey {
		return c
	}
	return c.ctx.Value(key)
}

// Cancel 取消上下文
func (c *Context) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
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

// Clone 克隆 Context（保留基础信息，新开 cancel 链路）
func (c *Context) Clone() *Context {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ctx, cancel := stdctx.WithCancel(stdctx.Background())
	clone := &Context{
		ctx:          ctx,
		cancel:       cancel,
		TraceID:      c.TraceID,
		StartTime:    time.Now(),
		UserID:       c.UserID,
		Username:     c.Username,
		UserRole:     c.UserRole,
		Method:       c.Method,
		Path:         c.Path,
		ClientIP:     c.ClientIP,
		UserAgent:    c.UserAgent,
		CustomFields: make(map[string]interface{}, len(c.CustomFields)),
		logger:       c.logger,
	}
	for k, v := range c.CustomFields {
		clone.CustomFields[k] = v
	}
	clone.logFields = make([]zap.Field, len(c.logFields))
	copy(clone.logFields, c.logFields)
	return clone
}
