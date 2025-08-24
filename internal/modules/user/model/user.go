package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Username        string         `json:"username" gorm:"not null;index"`
	Email           string         `json:"email" gorm:"not null;index"`
	Password        string         `json:"-" gorm:"not null"`
	PasswordVersion uint           `json:"-" gorm:"default:1"`
	Role            string         `json:"role" gorm:"default:'user'"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// 🔥 唯一索引将在数据库迁移中手动创建为条件索引，只对未删除的记录生效
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
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         Response `json:"user"`
}

// RefreshTokenRequest 刷新token请求结构
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse 刷新token响应结构
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
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

// AdminCreateUserRequest 管理员创建用户请求结构
type AdminCreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

// AdminUpdateUserRequest 管理员更新用户请求结构
type AdminUpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
	IsActive *bool  `json:"is_active" binding:"omitempty"`
}

// AdminResetPasswordRequest 管理员重置密码请求结构
type AdminResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
