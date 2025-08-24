package model

import (
	"time"

	storeModel "erp/internal/modules/store/model"
	supplierModel "erp/internal/modules/supplier/model"

	"gorm.io/gorm"
)

// SupplierInfo 供应商信息结构（用于API响应）
type SupplierInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
	IsActive bool   `json:"is_active"`
}

// StoreInfo 店铺信息结构（用于API响应）
type StoreInfo struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Remark     string `json:"remark"`
	IsActive   bool   `json:"is_active"`
	IsFeatured bool   `json:"is_featured"`
	SupplierID uint   `json:"supplier_id"`
}

// Sample 样品模型
type Sample struct {
	ID             uint                   `json:"id" gorm:"primaryKey"`
	ItemCode       string                 `json:"item_code" gorm:"not null;index;unique"`
	SupplierID     uint                   `json:"supplier_id" gorm:"not null;index"`
	StoreID        uint                   `json:"store_id" gorm:"not null;index"`
	HasLink        bool                   `json:"has_link" gorm:"default:false"`
	IsOffline      bool                   `json:"is_offline" gorm:"default:false"`
	CanModifyStock bool                   `json:"can_modify_stock" gorm:"default:true"`
	Supplier       supplierModel.Supplier `json:"supplier" gorm:"foreignKey:SupplierID"`
	Store          storeModel.Store       `json:"store" gorm:"foreignKey:StoreID"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      gorm.DeletedAt         `json:"-" gorm:"index"`
}

// CreateSampleRequest 创建样品请求结构
type CreateSampleRequest struct {
	ItemCode       string `json:"item_code" binding:"required,min=1,max=100"`
	SupplierID     uint   `json:"supplier_id" binding:"required"`
	StoreID        uint   `json:"store_id" binding:"required"`
	HasLink        *bool  `json:"has_link" binding:"omitempty"`
	IsOffline      *bool  `json:"is_offline" binding:"omitempty"`
	CanModifyStock *bool  `json:"can_modify_stock" binding:"omitempty"`
}

// UpdateSampleRequest 更新样品请求结构
type UpdateSampleRequest struct {
	ItemCode       string `json:"item_code" binding:"required,min=1,max=100"`
	SupplierID     uint   `json:"supplier_id" binding:"required"`
	StoreID        uint   `json:"store_id" binding:"required"`
	HasLink        *bool  `json:"has_link" binding:"omitempty"`
	IsOffline      *bool  `json:"is_offline" binding:"omitempty"`
	CanModifyStock *bool  `json:"can_modify_stock" binding:"omitempty"`
}

// SampleResponse 样品响应结构
type SampleResponse struct {
	ID             uint         `json:"id"`
	ItemCode       string       `json:"item_code"`
	SupplierID     uint         `json:"supplier_id"`
	StoreID        uint         `json:"store_id"`
	HasLink        bool         `json:"has_link"`
	IsOffline      bool         `json:"is_offline"`
	CanModifyStock bool         `json:"can_modify_stock"`
	Supplier       SupplierInfo `json:"supplier"`
	Store          StoreInfo    `json:"store"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// SampleListResponse 样品列表响应结构
type SampleListResponse struct {
	Samples    []SampleResponse `json:"samples"`
	Pagination Pagination       `json:"pagination"`
}

// Pagination 分页结构
type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// QuerySamplesRequest 查询样品请求参数
type QuerySamplesRequest struct {
	Page           int    `form:"page" binding:"omitempty,min=1"`
	Limit          int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search         string `form:"search"`
	SupplierID     *uint  `form:"supplier_id"`
	StoreID        *uint  `form:"store_id"`
	HasLink        *bool  `form:"has_link"`
	IsOffline      *bool  `form:"is_offline"`
	CanModifyStock *bool  `form:"can_modify_stock"`
	OrderBy        string `form:"order_by"`
}

// BatchUpdateSamplesRequest 批量更新样品请求结构
type BatchUpdateSamplesRequest struct {
	SampleIDs      []uint `json:"sample_ids" binding:"required,min=1"`
	HasLink        *bool  `json:"has_link" binding:"omitempty"`
	IsOffline      *bool  `json:"is_offline" binding:"omitempty"`
	CanModifyStock *bool  `json:"can_modify_stock" binding:"omitempty"`
}

// BatchUpdateSamplesResponse 批量更新样品响应结构
type BatchUpdateSamplesResponse struct {
	SuccessCount int    `json:"success_count"`
	FailedCount  int    `json:"failed_count"`
	TotalCount   int    `json:"total_count"`
	Message      string `json:"message"`
}
