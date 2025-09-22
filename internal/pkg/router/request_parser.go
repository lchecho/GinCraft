package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/constant"
	"github.com/liuchen/gin-craft/pkg/errors"
	"github.com/liuchen/gin-craft/pkg/response"
)

// ParseRequest 解析请求参数（泛型版本）
// 根据请求方法和Content-Type自动选择合适的绑定方式
func ParseRequest[T any](c *gin.Context) (*T, error) {
	var req T
	var err error

	// 根据请求方法和Content-Type选择绑定方式
	switch c.Request.Method {
	case "GET", "DELETE":
		// GET和DELETE请求使用Query参数
		err = c.ShouldBindQuery(&req)
	case "POST", "PUT", "PATCH":
		// 根据Content-Type选择绑定方式
		contentType := c.GetHeader("Content-Type")
		switch {
		case strings.Contains(contentType, "application/json"):
			// JSON格式
			err = c.ShouldBindJSON(&req)
		case strings.Contains(contentType, "application/x-www-form-urlencoded"):
			// Form表单格式
			err = c.ShouldBind(&req)
		case strings.Contains(contentType, "multipart/form-data"):
			// FormData格式（支持文件上传）
			err = c.ShouldBind(&req)
		default:
			// 默认尝试JSON绑定
			err = c.ShouldBindJSON(&req)
		}
	default:
		// 其他方法默认尝试JSON绑定
		err = c.ShouldBindJSON(&req)
	}

	// 参数绑定失败
	if err != nil {
		return nil, errors.New(constant.PARAM_ERROR, err.Error())
	}

	return &req, nil
}

// WithRequestHandler 包装带参数的处理函数
func WithRequestHandler[T any](handlerFunc func(*gin.Context, *T) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析请求参数
		req, err := ParseRequest[T](c)
		if err != nil {
			response.Error(c, err)
			return
		}

		// 执行业务逻辑
		data, err := handlerFunc(c, req)
		if err != nil {
			response.Error(c, err)
			return
		}

		// 返回成功响应
		response.Success(c, data)
	}
}
