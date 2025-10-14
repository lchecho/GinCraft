package errors

import (
	"errors"
	"fmt"

	"github.com/liuchen/gin-craft/internal/constant"
)

// AppError 应用错误结构
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("code: %d, message: %s, detail: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// GetCode 获取错误码
func (e *AppError) GetCode() int {
	return e.Code
}

// GetMessage 获取错误信息
func (e *AppError) GetMessage() string {
	return e.Message
}

// GetDetail 获取错误详情
func (e *AppError) GetDetail() string {
	return e.Detail
}

// New 创建新的应用错误
func New(code int, detail ...string) *AppError {
	message := constant.GetMsg(code)
	err := &AppError{
		Code:    code,
		Message: message,
	}
	if len(detail) > 0 && detail[0] != "" {
		err.Detail = detail[0]
	}
	return err
}

// Newf 创建新的应用错误（格式化详情）
func Newf(code int, format string, args ...interface{}) *AppError {
	message := constant.GetMsg(code)
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  fmt.Sprintf(format, args...),
	}
}

// IsAppError 判断是否为应用错误
func IsAppError(err error) bool {
	var appError *AppError
	ok := errors.As(err, &appError)
	return ok
}

// GetAppError 获取应用错误
func GetAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
