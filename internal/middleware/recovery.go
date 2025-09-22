package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Recovery 全局异常恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 绑定应用上下文（由 ContextMiddleware 提前注入）
				appCtx := MustGetContext(c)

				// 检查连接是否已断开
				// var brokenPipe bool
				// if ne, ok := err.(*net.OpError); ok {
				// 	if se, ok := ne.Err.(*os.SyscallError); ok {
				// 		if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
				// 			strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				// 			brokenPipe = true
				// 		}
				// 	}
				// }

				// 获取请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				// 获取堆栈信息
				stack := string(debug.Stack())
				// 获取错误信息
				errMsg := fmt.Sprintf("%v", err)

				// 添加日志字段
				appCtx.AddLogField("error", errMsg)
				appCtx.AddLogField("request", string(httpRequest))
				appCtx.AddLogField("stack", stack)
				// 记录日志
				appCtx.LogError("[Recovery from panic]")

				// 如果连接已断开，直接返回
				// if brokenPipe {
				// 	c.Error(err.(error))
				// 	c.Abort()
				// 	return
				// }

				// 返回错误响应
				c.Error(errors.WithStack(fmt.Errorf("%+v", err))).SetType(gin.ErrorTypePublic)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
