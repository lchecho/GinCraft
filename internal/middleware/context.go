package middleware

import (
	"github.com/gin-gonic/gin"
	pkgCtx "github.com/liuchen/gin-craft/internal/pkg/context"
	"github.com/liuchen/gin-craft/pkg/logger"
)

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
		c.Set(pkgCtx.CtxKey, appCtx)

		// 在响应头中添加追踪ID
		c.Header("X-Trace-ID", appCtx.GetTraceID())

		// 处理请求
		c.Next()

		// 取消应用Context
		// appCtx.Cancel()
	}
}
