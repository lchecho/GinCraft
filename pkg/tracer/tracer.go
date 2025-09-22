package tracer

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/liuchen/gin-craft/pkg/logger"
)

const (
	// TraceIDKey 跟踪ID的键名
	TraceIDKey = "trace_id"
)

// NewTraceID 生成新的跟踪ID
func NewTraceID() string {
	return uuid.New().String()
}

// ContextWithTraceID 在上下文中添加跟踪ID
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// TraceIDFromContext 从上下文中获取跟踪ID
func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// TraceIDFromGin 从Gin上下文中获取跟踪ID
func TraceIDFromGin(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if traceID, exists := c.Get(logger.TraceIDKey); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}

// SetTraceIDToGin 设置跟踪ID到Gin上下文
func SetTraceIDToGin(c *gin.Context, traceID string) {
	if c != nil {
		c.Set(logger.TraceIDKey, traceID)
		c.Header("X-Trace-ID", traceID)
	}
}
