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

// CreateProductWithAttributesRequest 创建产品（包含属性）请求结构
type CreateProductWithAttributesRequest struct {
	Name       string                       `json:"name" binding:"required,min=1,max=100"`
	CategoryID uint                         `json:"category_id" binding:"required,min=1"`
	Attributes []ProductAttributeValueInput `json:"attributes" binding:"required"` // 产品属性值
}

// ProductAttributeValueInput 产品属性值输入结构
type ProductAttributeValueInput struct {
	AttributeID uint        `json:"attribute_id" binding:"required"`
	Value       interface{} `json:"value" binding:"required"`
}

// UpdateProductRequest 更新产品请求结构
type UpdateProductRequest struct {
	Name       string `json:"name" binding:"omitempty,min=1,max=100"`
	CategoryID uint   `json:"category_id" binding:"omitempty,min=1"`
}

// UpdateProductWithAttributesRequest 更新产品（包含属性）请求结构
type UpdateProductWithAttributesRequest struct {
	Name       string                       `json:"name" binding:"omitempty,min=1,max=100"`
	CategoryID uint                         `json:"category_id" binding:"omitempty,min=1"`
	Attributes []ProductAttributeValueInput `json:"attributes" binding:"omitempty"` // 产品属性值
}

// ProductResponse 产品响应结构
type ProductResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ProductWithAttributesResponse 产品（包含属性）响应结构
type ProductWithAttributesResponse struct {
	ID         uint                       `json:"id"`
	Name       string                     `json:"name"`
	CategoryID uint                       `json:"category_id"`
	Attributes []ProductAttributeResponse `json:"attributes"` // 产品属性值
	CreatedAt  time.Time                  `json:"created_at"`
	UpdatedAt  time.Time                  `json:"updated_at"`
}

// ProductAttributeResponse 产品属性响应结构
type ProductAttributeResponse struct {
	AttributeID   uint        `json:"attribute_id"`
	AttributeName string      `json:"attribute_name"`
	DisplayName   string      `json:"display_name"`
	AttributeType string      `json:"attribute_type"`
	Value         interface{} `json:"value"`
	IsRequired    bool        `json:"is_required"`
	IsInherited   bool        `json:"is_inherited"`   // 是否为继承属性
	InheritedFrom *uint       `json:"inherited_from"` // 继承来源分类ID
}

// ProductListResponse 产品列表响应结构
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Pagination Pagination        `json:"pagination"`
}

// ProductWithAttributesListResponse 产品（包含属性）列表响应结构
type ProductWithAttributesListResponse struct {
	Products   []ProductWithAttributesResponse `json:"products"`
	Pagination Pagination                      `json:"pagination"`
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

// CategoryAttributeTemplateResponse 分类属性模板响应（用于产品创建表单）
type CategoryAttributeTemplateResponse struct {
	CategoryID uint                                    `json:"category_id"`
	Attributes []CategoryAttributeTemplateItemResponse `json:"attributes"`
}

// CategoryAttributeTemplateItemResponse 分类属性模板项响应
type CategoryAttributeTemplateItemResponse struct {
	AttributeID   uint        `json:"attribute_id"`
	Name          string      `json:"name"`
	DisplayName   string      `json:"display_name"`
	Type          string      `json:"type"`
	Unit          string      `json:"unit"`
	IsRequired    bool        `json:"is_required"`
	DefaultValue  string      `json:"default_value"`
	Options       interface{} `json:"options"`    // 选项配置
	Validation    interface{} `json:"validation"` // 验证规则
	Sort          int         `json:"sort"`
	IsInherited   bool        `json:"is_inherited"`   // 是否为继承属性
	InheritedFrom *uint       `json:"inherited_from"` // 继承来源分类ID
}

// ValidateProductAttributesRequest 验证产品属性请求结构
type ValidateProductAttributesRequest struct {
	CategoryID uint                         `json:"category_id" binding:"required"`
	Attributes []ProductAttributeValueInput `json:"attributes" binding:"required"`
}

// ValidationResult 验证结果
type ValidationResult struct {
	IsValid bool                       `json:"is_valid"`
	Errors  []AttributeValidationError `json:"errors,omitempty"`
}

// AttributeValidationError 属性验证错误
type AttributeValidationError struct {
	AttributeID uint   `json:"attribute_id"`
	Field       string `json:"field"`
	Message     string `json:"message"`
}
