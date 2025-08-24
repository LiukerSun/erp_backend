package model

import (
	"time"

	storeModel "erp/internal/modules/store/model"

	"gorm.io/gorm"
)

// StoreInfo 店铺信息结构（用于API响应）
type StoreInfo struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Remark     string    `json:"remark"`
	IsActive   bool      `json:"is_active"`
	SupplierID uint      `json:"supplier_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Supplier 供应商模型
type Supplier struct {
	ID        uint               `json:"id" gorm:"primaryKey"`
	Name      string             `json:"name" gorm:"not null;index"`
	Remark    string             `json:"remark" gorm:"type:text"`
	IsActive  bool               `json:"is_active" gorm:"default:true"`
	Stores    []storeModel.Store `json:"stores" gorm:"foreignKey:SupplierID"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `json:"-" gorm:"index"`
}

// CreateSupplierRequest 创建供应商请求结构
type CreateSupplierRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=100"`
	Remark string `json:"remark" binding:"max=500"`
}

// UpdateSupplierRequest 更新供应商请求结构
type UpdateSupplierRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Remark   string `json:"remark" binding:"max=500"`
	IsActive *bool  `json:"is_active" binding:"omitempty"`
}

// SupplierResponse 供应商响应结构
type SupplierResponse struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Remark    string      `json:"remark"`
	IsActive  bool        `json:"is_active"`
	Stores    []StoreInfo `json:"stores"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// SupplierListResponse 供应商列表响应结构
type SupplierListResponse struct {
	Suppliers  []SupplierResponse `json:"suppliers"`
	Pagination Pagination         `json:"pagination"`
}

// Pagination 分页结构
type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// QuerySuppliersRequest 查询供应商请求参数
type QuerySuppliersRequest struct {
	Page          int    `form:"page" binding:"omitempty,min=1"`
	Limit         int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search        string `form:"search"`
	IncludeStores bool   `form:"include_stores"`
	IsActive      *bool  `form:"is_active"` // 筛选活跃状态：true=活跃，false=不活跃，nil=所有
	OrderBy       string `form:"order_by"`
}
