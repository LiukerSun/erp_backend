package model

import (
	"time"

	"gorm.io/gorm"
)

// Store 店铺模型
type Store struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	Name                  string         `json:"name" gorm:"not null;index"`
	Remark                string         `json:"remark" gorm:"type:text"`
	IsActive              bool           `json:"is_active" gorm:"default:true"`
	IsFeatured            bool           `json:"is_featured" gorm:"default:false"`
	DefaultCommissionRate int            `json:"default_commission_rate" gorm:"default:0"` // 默认佣金率 0-100 1 表示1%
	SupplierID            uint           `json:"supplier_id" gorm:"not null;index"`
	Supplier              Supplier       `json:"supplier" gorm:"foreignKey:SupplierID"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}

// Supplier 供应商模型（用于关联）
type Supplier struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"not null"`
	Remark string `json:"remark" gorm:"type:text"`
}

// CreateStoreRequest 创建店铺请求结构
type CreateStoreRequest struct {
	Name                  string `json:"name" binding:"required,min=1,max=100"`
	Remark                string `json:"remark" binding:"max=500"`
	SupplierID            uint   `json:"supplier_id" binding:"required"`
	IsFeatured            *bool  `json:"is_featured" binding:"omitempty"`
	DefaultCommissionRate int    `json:"default_commission_rate" binding:"omitempty"` // 默认佣金率 0-100 1 表示1%
}

// UpdateStoreRequest 更新店铺请求结构
type UpdateStoreRequest struct {
	Name                  string `json:"name" binding:"required,min=1,max=100"`
	Remark                string `json:"remark" binding:"max=500"`
	SupplierID            uint   `json:"supplier_id" binding:"required"`
	IsActive              *bool  `json:"is_active" binding:"omitempty"`
	IsFeatured            *bool  `json:"is_featured" binding:"omitempty"`
	DefaultCommissionRate int    `json:"default_commission_rate" binding:"omitempty"` // 默认佣金率 0-100 1 表示1%
}

// StoreResponse 店铺响应结构
type StoreResponse struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Remark                string    `json:"remark"`
	IsActive              bool      `json:"is_active"`
	IsFeatured            bool      `json:"is_featured"`
	DefaultCommissionRate int       `json:"default_commission_rate" gorm:"default:0"` // 默认佣金率 0-100 1 表示1%
	SupplierID            uint      `json:"supplier_id"`
	Supplier              Supplier  `json:"supplier"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// StoreListResponse 店铺列表响应结构
type StoreListResponse struct {
	Stores     []StoreResponse `json:"stores"`
	Pagination Pagination      `json:"pagination"`
}

// Pagination 分页结构
type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
