package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	pkgCtx "github.com/liuchen/gin-craft/internal/pkg/context"
	"github.com/liuchen/gin-craft/pkg/logger"
)

// ContextKey 用于在gin.Context中存储自定义Context的键
const ContextKey = "custom_context"

// ContextMiddleware 自定义Context中间件
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建应用Context
		appCtx := pkgCtx.New(c.Request.Context())

		// 设置请求信息
		appCtx.SetRequestInfo(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		// 设置日志记录器
		appCtx.SetLogger(logger.Log)

		// 记录trace_id
		c.Set(logger.TraceIDKey, appCtx.GetTraceID())

		// 将应用Context存储到gin.Context中
		c.Set(ContextKey, appCtx)

		// 在响应头中添加追踪ID
		c.Header("X-Trace-ID", appCtx.GetTraceID())

		// 处理请求
		c.Next()

		// 取消应用Context
		// appCtx.Cancel()
	}
}

// GetContext 从gin.Context中获取应用Context
func GetContext(c context.Context) *pkgCtx.Context {
	value := c.Value(ContextKey)
	if appCtx, ok := value.(*pkgCtx.Context); ok {
		return appCtx
	}

	return nil
}

// MustGetContext 从gin.Context中获取应用Context，如果不存在则panic
func MustGetContext(c context.Context) *pkgCtx.Context {
	appCtx := GetContext(c)
	if appCtx == nil {
		panic("App context not found in gin.Context")
	}
	return appCtx
}
