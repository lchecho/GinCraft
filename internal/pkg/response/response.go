package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/constant"
	"github.com/liuchen/gin-craft/internal/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Detail interface{} `json:"detail,omitempty"`
}

// httpStatusByCode 业务错误码 → HTTP 状态码映射
var httpStatusByCode = map[int]int{
	constant.Unauthorized:    http.StatusUnauthorized,
	constant.Forbidden:       http.StatusForbidden,
	constant.NotFound:        http.StatusNotFound,
	constant.MethodNotAllow:  http.StatusMethodNotAllowed,
	constant.TooManyRequests: http.StatusTooManyRequests,
	constant.Timeout:         http.StatusGatewayTimeout,
	constant.ParamError:      http.StatusBadRequest,
	constant.SystemError:     http.StatusInternalServerError,
	constant.DBError:         http.StatusInternalServerError,
}

func httpStatusOf(code int) int {
	if s, ok := httpStatusByCode[code]; ok {
		return s
	}
	return http.StatusOK
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: constant.Success,
		Msg:  constant.GetMsg(constant.Success),
		Data: data,
	})
}

// Fail 失败响应（基于业务错误码自动选择 HTTP 状态码）
func Fail(c *gin.Context, code int, data interface{}) {
	c.JSON(httpStatusOf(code), Response{
		Code: code,
		Msg:  constant.GetMsg(code),
		Data: data,
	})
}

// FailWithMsg 自定义消息的失败响应
func FailWithMsg(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(httpStatusOf(code), Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// FailWithDetail 带详情的失败响应
func FailWithDetail(c *gin.Context, code int, detail string, data interface{}) {
	c.JSON(httpStatusOf(code), Response{
		Code:   code,
		Msg:    constant.GetMsg(code),
		Data:   data,
		Detail: detail,
	})
}

// Error 错误响应（支持 AppError，自动选择 HTTP 状态码）
func Error(c *gin.Context, err error) {
	if err == nil {
		Success(c, nil)
		return
	}
	if appErr, ok := errors.GetAppError(err); ok {
		resp := Response{
			Code: appErr.GetCode(),
			Msg:  appErr.GetMessage(),
		}
		if d := appErr.GetDetail(); d != "" {
			resp.Detail = d
		}
		c.JSON(httpStatusOf(appErr.GetCode()), resp)
		return
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:   constant.SystemError,
		Msg:    constant.GetMsg(constant.SystemError),
		Detail: err.Error(),
	})
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
