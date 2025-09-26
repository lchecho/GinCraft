package middleware

import (
	"bytes"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装gin.ResponseWriter以捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定应用上下文（由 ContextMiddleware 提前注入）
		appCtx := MustGetContext(c)
		// 请求路径
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		// 请求方法
		method := c.Request.Method
		// 请求IP
		clientIP := c.ClientIP()
		// User-Agent
		userAgent := c.Request.UserAgent()

		// 获取请求体（排除文件上传）
		var requestBody []byte
		contentType := c.Request.Header.Get("Content-Type")
		if c.Request.Body != nil && shouldLogRequestBody(contentType) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 构建日志字段
		appCtx.AddLogField("method", method)
		appCtx.AddLogField("path", path)
		appCtx.AddLogField("ip", clientIP)
		appCtx.AddLogField("user_agent", userAgent)
		// 添加请求体（如果需要记录）
		if len(requestBody) > 0 {
			appCtx.AddLogField("request_body", string(requestBody))
		} else if isFileUpload(contentType) {
			// 文件上传请求，只记录类型不记录内容
			appCtx.AddLogField("request_type", "file_upload")
		}

		// 包装ResponseWriter以捕获响应体
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 状态码
		statusCode := c.Writer.Status()
		// 响应大小
		responseSize := c.Writer.Size()
		// 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		appCtx.AddLogField("status", statusCode)
		appCtx.AddLogField("response_size", responseSize)

		// 添加响应体（如果需要记录且不是文件下载等）
		if shouldLogResponseBody(c.Writer.Header().Get("Content-Type")) && w.body.Len() > 0 && w.body.Len() < 1024 {
			appCtx.AddLogField("response_body", w.body.String())
		}

		// 添加错误信息
		if errorMessage != "" {
			appCtx.AddLogField("error", errorMessage)
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			appCtx.LogError("HTTP Request")
		} else if statusCode >= 400 {
			appCtx.LogWarn("HTTP Request")
		} else {
			appCtx.LogInfo("HTTP Request")
		}
	}
}

// shouldLogRequestBody 判断是否应该记录请求体
func shouldLogRequestBody(contentType string) bool {
	if contentType == "" {
		return false
	}

	// 只记录JSON和普通表单数据，不记录文件上传
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "application/x-www-form-urlencoded")
}

// isFileUpload 判断是否为文件上传请求
func isFileUpload(contentType string) bool {
	return strings.Contains(contentType, "multipart/form-data")
}

// shouldLogResponseBody 判断是否应该记录响应体
func shouldLogResponseBody(contentType string) bool {
	if contentType == "" {
		return true // 默认记录
	}

	// 不记录文件下载、图片等二进制内容
	return !strings.Contains(contentType, "application/octet-stream") &&
		!strings.Contains(contentType, "image/") &&
		!strings.Contains(contentType, "video/") &&
		!strings.Contains(contentType, "audio/")
}
