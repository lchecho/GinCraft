package user

import (
	"github.com/liuchen/gin-craft/internal/dto"
)

// RegisterRequest 用户注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20" example:"john_doe"` // 用户名，3-20个字符
	Password string `json:"password" binding:"required,min=6,max=20" example:"123456"`   // 密码，6-20个字符
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`   // 邮箱地址
}

// LoginRequest 用户登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"` // 用户名
	Password string `json:"password" binding:"required" example:"123456"`   // 密码
}

// UpdateRequest 用户更新请求参数
type UpdateRequest struct {
	ID       uint   `json:"id" binding:"omitempty"`
	Username string `json:"username" binding:"omitempty,min=3,max=20" example:"john_doe"` // 用户名，3-20个字符
	Email    string `json:"email" binding:"omitempty,email" example:"john@example.com"`   // 邮箱地址
}

// PasswordUpdateRequest 用户密码更新请求参数
type PasswordUpdateRequest struct {
	ID       uint   `json:"id" binding:"omitempty"`
	Password string `json:"password" binding:"required,min=6,max=20" example:"newpassword123"` // 新密码，6-20个字符
}

// ListRequest 用户列表请求参数
type ListRequest struct {
	dto.Pagination
	Username string `form:"username" json:"username" example:"john"` // 用户名筛选
	Email    string `form:"email" json:"email" example:"john@"`      // 邮箱筛选
}

// InfoRequest 获取用户信息请求参数
type InfoRequest struct {
	ID uint `json:"id" binding:"required"`
}
