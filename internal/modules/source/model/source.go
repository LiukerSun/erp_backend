package model

import (
	"time"
)

// Source 货源模型
// @Description 货源信息
type Source struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null" example:"Apple官方旗舰店"`          // 货源名称
	Code      string     `json:"code" gorm:"type:varchar(50);uniqueIndex;not null" example:"APPLE001"` // 货源编码
	Status    int        `json:"status" gorm:"default:1" example:"1"`                                  // 状态：1-启用，0-禁用
	Remark    string     `json:"remark" gorm:"type:text" example:"优质货源"`                               // 备注
}
