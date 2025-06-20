package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Username        string         `json:"username" gorm:"uniqueIndex;not null"`
	Email           string         `json:"email" gorm:"uniqueIndex;not null"`
	Password        string         `json:"-" gorm:"not null"`  // 密码不返回给前端
	PasswordVersion uint           `json:"-" gorm:"default:1"` // 密码版本，用于使旧token失效
	Role            string         `json:"role" gorm:"default:'user'"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Response 用户响应结构
type Response struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateProfileRequest 更新资料请求结构
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"email"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string   `json:"token"`
	User  Response `json:"user"`
}

// UserListResponse 用户列表响应结构
type UserListResponse struct {
	Users      []Response `json:"users"`
	Pagination Pagination `json:"pagination"`
}

// Pagination 分页结构
type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
