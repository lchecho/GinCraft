package user

import (
	"github.com/liuchen/gin-craft/internal/dto"
	"time"
)

// User 用户响应参数
type User struct {
	ID        uint      `json:"id" example:"1"`                            // 用户ID
	Username  string    `json:"username" example:"john_doe"`               // 用户名
	Email     string    `json:"email" example:"john@example.com"`          // 邮箱地址
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"` // 更新时间
}

// LoginResponse 用户登录响应参数
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 访问令牌
}

// ListResponse 用户列表响应参数
type ListResponse struct {
	List []User `json:"list"`
	dto.Pagination
}
