package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/constant"
	"github.com/liuchen/gin-craft/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Detail string      `json:"detail,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: constant.Success,
		Msg:  constant.GetMsg(constant.Success),
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  constant.GetMsg(code),
		Data: data,
	})
}

// FailWithMsg 自定义消息的失败响应
func FailWithMsg(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// FailWithDetail 带详情的失败响应
func FailWithDetail(c *gin.Context, code int, detail string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:   code,
		Msg:    constant.GetMsg(code),
		Data:   data,
		Detail: detail,
	})
}

// Error 错误响应（支持AppError）
func Error(c *gin.Context, err error) {
	if appErr, ok := errors.GetAppError(err); ok {
		// 处理应用错误
		resp := Response{
			Code: appErr.GetCode(),
			Msg:  appErr.GetMessage(),
			Data: nil,
		}
		if detail := appErr.GetDetail(); detail != "" {
			resp.Detail = detail
		}
		c.JSON(http.StatusOK, resp)
	} else {
		// 处理普通错误
		c.JSON(http.StatusOK, Response{
			Code:   constant.SystemError,
			Msg:    constant.GetMsg(constant.SystemError),
			Data:   nil,
			Detail: err.Error(),
		})
	}
}

// ParamError 参数错误响应
func ParamError(c *gin.Context) {
	Fail(c, constant.ParamError, nil)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context) {
	Fail(c, constant.SystemError, nil)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context) {
	Fail(c, constant.Unauthorized, nil)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context) {
	Fail(c, constant.Forbidden, nil)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context) {
	Fail(c, constant.NotFound, nil)
}

// TooManyRequests 请求过多响应
func TooManyRequests(c *gin.Context) {
	Fail(c, constant.TooManyRequests, nil)
}
