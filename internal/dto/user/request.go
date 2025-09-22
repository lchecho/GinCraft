package user

import "github.com/liuchen/gin-craft/internal/dto/common"

// UserRegisterRequest 用户注册请求参数
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20" example:"john_doe"` // 用户名，3-20个字符
	Password string `json:"password" binding:"required,min=6,max=20" example:"123456"`   // 密码，6-20个字符
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`   // 邮箱地址
}

// UserLoginRequest 用户登录请求参数
type UserLoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"` // 用户名
	Password string `json:"password" binding:"required" example:"123456"`   // 密码
}

// UserUpdateRequest 用户更新请求参数
type UserUpdateRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20" example:"john_doe"` // 用户名，3-20个字符
	Email    string `json:"email" binding:"omitempty,email" example:"john@example.com"`   // 邮箱地址
}

// UserPasswordUpdateRequest 用户密码更新请求参数
type UserPasswordUpdateRequest struct {
	Password string `json:"password" binding:"required,min=6,max=20" example:"newpassword123"` // 新密码，6-20个字符
}

// UserListRequest 用户列表请求参数
type UserListRequest struct {
	common.PaginationRequest
	Username string `form:"username" json:"username" example:"john"` // 用户名筛选
	Email    string `form:"email" json:"email" example:"john@"`      // 邮箱筛选
}

// UserInfoRequest 获取用户信息请求参数（空结构体，因为信息从token中获取）
type UserInfoRequest struct{}
