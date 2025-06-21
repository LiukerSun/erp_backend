package model

import (
	"erp/internal/modules/product/model"
	"time"
)

// Tag 标签模型
// @Description 标签信息
type Tag struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   *time.Time      `json:"deletedAt,omitempty" gorm:"index"`
	Name        string          `json:"name" gorm:"type:varchar(50);uniqueIndex;not null" example:"热销"` // 标签名称
	Description string          `json:"description" gorm:"type:varchar(200)" example:"热销商品标签"`          // 标签描述
	Color       string          `json:"color" gorm:"type:varchar(7)" example:"#FF6B6B"`                 // 标签颜色
	IsEnabled   bool            `json:"is_enabled" gorm:"default:true" example:"true"`                  // 是否启用
	Products    []model.Product `json:"products" gorm:"many2many:product_tags;"`                        // 关联的商品
}

// ProductTag 商品和标签的多对多关联表
type ProductTag struct {
	ProductID uint      `gorm:"primaryKey;column:product_id"`
	TagID     uint      `gorm:"primaryKey;column:tag_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
