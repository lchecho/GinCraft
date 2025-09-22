package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"username"`
	Password  string         `gorm:"type:varchar(100);not null" json:"-"` // 不返回密码
	Email     string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
