package model

import (
	"time"

	"gorm.io/gorm"
)

// Product 产品模型
type Product struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"not null"`        // 产品名（可重复）
	CategoryID uint           `json:"category_id" gorm:"not null"` // 产品分类（数字）
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateProductRequest 创建产品请求结构
type CreateProductRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=100"`
	CategoryID uint   `json:"category_id" binding:"required,min=1"`
}

// UpdateProductRequest 更新产品请求结构
type UpdateProductRequest struct {
	Name       string `json:"name" binding:"omitempty,min=1,max=100"`
	CategoryID uint   `json:"category_id" binding:"omitempty,min=1"`
}

// ProductResponse 产品响应结构
type ProductResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ProductListResponse 产品列表响应结构
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Pagination Pagination        `json:"pagination"`
}

// Pagination 分页结构
type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// ProductQueryRequest 产品查询请求结构
type ProductQueryRequest struct {
	Name       string `form:"name"`        // 按名称模糊搜索
	CategoryID uint   `form:"category_id"` // 按分类筛选
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
}
