package middleware

import (
	stdctx "context"

	"github.com/gin-gonic/gin"
	pkgCtx "github.com/liuchen/gin-craft/internal/pkg/context"
	"github.com/liuchen/gin-craft/pkg/logger"
)

// ContextMiddleware 构造应用 Context 并注入到 gin + Request.Context
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许上游透传 TraceID
		traceID := logger.TraceIDFromHeaders(c.Request.Header)
		appCtx := pkgCtx.NewWithTraceID(c.Request.Context(), traceID)
		appCtx.SetRequestInfo(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
		appCtx.SetLogger(logger.Log)

		// 1. traceID 放入 gin.Keys（string key，便于老式 c.MustGet 读取）
		c.Set(logger.TraceIDKey, appCtx.GetTraceID())

		// 2. appCtx 放入 Request.Context（类型化 key，可跨 goroutine / 下游 ctx 传递）
		ctx := stdctx.WithValue(c.Request.Context(), pkgCtx.CtxKey, appCtx)
		c.Request = c.Request.WithContext(ctx)

		// 响应头暴露 TraceID，便于链路追踪
		c.Header("X-Trace-ID", appCtx.GetTraceID())

		c.Next()
	}
}
