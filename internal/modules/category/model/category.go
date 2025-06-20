package model

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;index"`    // 分类名称
	Description string         `json:"description"`                   // 分类描述
	ParentID    *uint          `json:"parent_id" gorm:"index"`        // 父分类ID，为空表示根分类
	Level       int            `json:"level" gorm:"default:1"`        // 分类层级
	Sort        int            `json:"sort" gorm:"default:0"`         // 排序字段
	IsActive    bool           `json:"is_active" gorm:"default:true"` // 是否启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Parent   *Category   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`   // 父分类
	Children []Category  `json:"children,omitempty" gorm:"foreignKey:ParentID"` // 子分类
	Products interface{} `json:"products,omitempty" gorm:"-"`                   // 关联的产品（使用interface{}避免循环依赖）
}

// CreateCategoryRequest 创建分类请求结构
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	ParentID    *uint  `json:"parent_id" binding:"omitempty"`
	Sort        int    `json:"sort" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateCategoryRequest 更新分类请求结构
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	ParentID    *uint  `json:"parent_id" binding:"omitempty"`
	Sort        int    `json:"sort" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// CategoryResponse 分类响应结构
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	Level       int       `json:"level"`
	Sort        int       `json:"sort"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryTreeResponse 分类树响应结构
type CategoryTreeResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ParentID    *uint                   `json:"parent_id"`
	Level       int                     `json:"level"`
	Sort        int                     `json:"sort"`
	IsActive    bool                    `json:"is_active"`
	Children    []*CategoryTreeResponse `json:"children,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

// CategoryListResponse 分类列表响应结构
type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Pagination Pagination         `json:"pagination"`
}

// CategoryTreeListResponse 分类树列表响应结构
type CategoryTreeListResponse struct {
	Categories []*CategoryTreeResponse `json:"categories"`
}

// Pagination 分页结构
type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// CategoryQueryRequest 分类查询请求结构
type CategoryQueryRequest struct {
	Name     string `form:"name"`      // 按名称模糊搜索
	ParentID *uint  `form:"parent_id"` // 按父分类筛选
	IsActive *bool  `form:"is_active"` // 按状态筛选
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

// MoveCategoryRequest 移动分类请求结构
type MoveCategoryRequest struct {
	ParentID *uint `json:"parent_id" binding:"omitempty"`
}

// CategoryPathResponse 分类路径响应结构
type CategoryPathResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// CategoryWithPathResponse 带路径的分类响应结构
type CategoryWithPathResponse struct {
	CategoryResponse
	Path []CategoryPathResponse `json:"path"` // 从根到当前分类的路径
}
