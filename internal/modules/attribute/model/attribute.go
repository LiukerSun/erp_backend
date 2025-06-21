package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// AttributeType 属性类型枚举
type AttributeType string

const (
	AttributeTypeText        AttributeType = "text"         // 文本类型
	AttributeTypeNumber      AttributeType = "number"       // 数字类型
	AttributeTypeSelect      AttributeType = "select"       // 单选类型
	AttributeTypeMultiSelect AttributeType = "multi_select" // 多选类型
	AttributeTypeBoolean     AttributeType = "boolean"      // 布尔类型
	AttributeTypeDate        AttributeType = "date"         // 日期类型
	AttributeTypeDateTime    AttributeType = "datetime"     // 日期时间类型
	AttributeTypeURL         AttributeType = "url"          // URL类型
	AttributeTypeEmail       AttributeType = "email"        // 邮箱类型
	AttributeTypeColor       AttributeType = "color"        // 颜色类型
	AttributeTypeCurrency    AttributeType = "currency"     // 货币类型
)

// AttributeValueType 属性值存储类型
type AttributeValueType string

const (
	ValueTypeText   AttributeValueType = "text"
	ValueTypeNumber AttributeValueType = "number"
	ValueTypeJSON   AttributeValueType = "json"
	ValueTypeBool   AttributeValueType = "bool"
	ValueTypeDate   AttributeValueType = "date"
)

// Attribute 属性模型
type Attribute struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null;index"`       // 🔥 移除unique，改为普通index
	DisplayName  string         `json:"display_name" gorm:"not null"`     // 显示名称
	Description  string         `json:"description"`                      // 属性描述
	Type         AttributeType  `json:"type" gorm:"not null"`             // 属性类型
	Unit         string         `json:"unit"`                             // 单位（如：kg, cm, 元等）
	IsRequired   bool           `json:"is_required" gorm:"default:false"` // 是否必填
	DefaultValue string         `json:"default_value"`                    // 默认值
	Options      string         `json:"options" gorm:"type:text"`         // 选项配置（JSON格式）
	Validation   string         `json:"validation" gorm:"type:text"`      // 验证规则（JSON格式）
	Sort         int            `json:"sort" gorm:"default:0"`            // 排序
	IsActive     bool           `json:"is_active" gorm:"default:true"`    // 是否启用
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Values            []AttributeValue    `json:"values,omitempty" gorm:"foreignKey:AttributeID"`              // 属性值
	CategoryAttribute []CategoryAttribute `json:"category_attributes,omitempty" gorm:"foreignKey:AttributeID"` // 分类属性关联

	// 🔥 唯一索引将在数据库迁移中手动创建为条件索引，只对未删除的记录生效
}

// IsMultiple 判断属性是否支持多值
func (a *Attribute) IsMultiple() bool {
	return a.Type == AttributeTypeMultiSelect
}

// SupportsOptions 判断属性是否支持选项配置
func (a *Attribute) SupportsOptions() bool {
	return a.Type == AttributeTypeSelect || a.Type == AttributeTypeMultiSelect
}

// CategoryAttribute 分类属性关联模型（多对多中间表）
type CategoryAttribute struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  uint           `json:"category_id" gorm:"not null;index"`  // 分类ID
	AttributeID uint           `json:"attribute_id" gorm:"not null;index"` // 属性ID
	IsRequired  bool           `json:"is_required" gorm:"default:false"`   // 在此分类中是否必填
	Sort        int            `json:"sort" gorm:"default:0"`              // 在分类中的排序
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Attribute Attribute `json:"attribute,omitempty" gorm:"foreignKey:AttributeID"`

	// 🔥 修复软删除唯一索引问题：使用条件唯一索引，只对未删除的记录生效
	// 这个唯一索引将在数据库迁移中手动创建，而不是依赖GORM的自动创建
}

// AttributeValue 属性值模型
type AttributeValue struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	AttributeID uint               `json:"attribute_id" gorm:"not null;index"` // 属性ID
	EntityType  string             `json:"entity_type" gorm:"not null;index"`  // 实体类型（product, category等）
	EntityID    uint               `json:"entity_id" gorm:"not null;index"`    // 实体ID
	ValueType   AttributeValueType `json:"value_type" gorm:"not null"`         // 值类型
	TextValue   string             `json:"text_value"`                         // 文本值
	NumberValue *float64           `json:"number_value"`                       // 数字值
	BoolValue   *bool              `json:"bool_value"`                         // 布尔值
	DateValue   *time.Time         `json:"date_value"`                         // 日期值
	JSONValue   string             `json:"json_value" gorm:"type:text"`        // JSON值（用于复杂类型）
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `json:"-" gorm:"index"`

	// 关联关系
	Attribute Attribute `json:"attribute,omitempty" gorm:"foreignKey:AttributeID"`

	// 🔥 修复软删除唯一索引问题：使用条件唯一索引，只对未删除的记录生效
	// 这个唯一索引将在数据库迁移中手动创建，而不是依赖GORM的自动创建
}

// AttributeOption 属性选项结构（用于select和multi_select类型）
type AttributeOption struct {
	Value       string                 `json:"value"`
	Label       string                 `json:"label"`
	Color       string                 `json:"color,omitempty"`
	Description string                 `json:"description,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}

// ValidationRule 验证规则结构
type ValidationRule struct {
	MinLength *int                   `json:"min_length,omitempty"` // 最小长度
	MaxLength *int                   `json:"max_length,omitempty"` // 最大长度
	Min       *float64               `json:"min,omitempty"`        // 最小值
	Max       *float64               `json:"max,omitempty"`        // 最大值
	Pattern   string                 `json:"pattern,omitempty"`    // 正则表达式
	Required  bool                   `json:"required"`             // 是否必填
	Custom    map[string]interface{} `json:"custom,omitempty"`     // 自定义规则
}

// GetOptions 获取属性选项
func (a *Attribute) GetOptions() ([]AttributeOption, error) {
	if a.Options == "" {
		return []AttributeOption{}, nil
	}

	var options []AttributeOption
	err := json.Unmarshal([]byte(a.Options), &options)
	return options, err
}

// SetOptions 设置属性选项
func (a *Attribute) SetOptions(options []AttributeOption) error {
	data, err := json.Marshal(options)
	if err != nil {
		return err
	}
	a.Options = string(data)
	return nil
}

// GetValidation 获取验证规则
func (a *Attribute) GetValidation() (ValidationRule, error) {
	if a.Validation == "" {
		return ValidationRule{}, nil
	}

	var validation ValidationRule
	err := json.Unmarshal([]byte(a.Validation), &validation)
	return validation, err
}

// SetValidation 设置验证规则
func (a *Attribute) SetValidation(validation ValidationRule) error {
	data, err := json.Marshal(validation)
	if err != nil {
		return err
	}
	a.Validation = string(data)
	return nil
}

// GetValue 获取属性值的实际值
func (av *AttributeValue) GetValue() interface{} {
	switch av.ValueType {
	case ValueTypeText:
		return av.TextValue
	case ValueTypeNumber:
		return av.NumberValue
	case ValueTypeBool:
		return av.BoolValue
	case ValueTypeDate:
		return av.DateValue
	case ValueTypeJSON:
		var value interface{}
		json.Unmarshal([]byte(av.JSONValue), &value)
		return value
	default:
		return nil
	}
}

// SetValue 设置属性值
func (av *AttributeValue) SetValue(value interface{}) error {
	switch v := value.(type) {
	case string:
		av.ValueType = ValueTypeText
		av.TextValue = v
	case int, int32, int64, float32, float64:
		av.ValueType = ValueTypeNumber
		if floatVal, ok := value.(float64); ok {
			av.NumberValue = &floatVal
		} else {
			floatVal := float64(v.(int))
			av.NumberValue = &floatVal
		}
	case bool:
		av.ValueType = ValueTypeBool
		av.BoolValue = &v
	case time.Time:
		av.ValueType = ValueTypeDate
		av.DateValue = &v
	default:
		av.ValueType = ValueTypeJSON
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		av.JSONValue = string(data)
	}
	return nil
}

// CreateAttributeRequest 创建属性请求结构
type CreateAttributeRequest struct {
	Name         string            `json:"name" binding:"required,min=1,max=100"`
	DisplayName  string            `json:"display_name" binding:"required,min=1,max=100"`
	Description  string            `json:"description" binding:"omitempty,max=500"`
	Type         AttributeType     `json:"type" binding:"required,oneof=text number select multi_select boolean date datetime url email color currency"`
	Unit         string            `json:"unit" binding:"omitempty,max=20"`
	IsRequired   *bool             `json:"is_required" binding:"omitempty"`
	DefaultValue string            `json:"default_value" binding:"omitempty"`
	Options      []AttributeOption `json:"options" binding:"omitempty"`
	Validation   ValidationRule    `json:"validation" binding:"omitempty"`
	Sort         int               `json:"sort" binding:"omitempty,min=0"`
	IsActive     *bool             `json:"is_active" binding:"omitempty"`
}

// UpdateAttributeRequest 更新属性请求结构
type UpdateAttributeRequest struct {
	Name         string            `json:"name" binding:"omitempty,min=1,max=100"`
	DisplayName  string            `json:"display_name" binding:"omitempty,min=1,max=100"`
	Description  string            `json:"description" binding:"omitempty,max=500"`
	Type         AttributeType     `json:"type" binding:"omitempty,oneof=text number select multi_select boolean date datetime url email color currency"`
	Unit         string            `json:"unit" binding:"omitempty,max=20"`
	IsRequired   *bool             `json:"is_required" binding:"omitempty"`
	DefaultValue string            `json:"default_value" binding:"omitempty"`
	Options      []AttributeOption `json:"options" binding:"omitempty"`
	Validation   ValidationRule    `json:"validation" binding:"omitempty"`
	Sort         int               `json:"sort" binding:"omitempty,min=0"`
	IsActive     *bool             `json:"is_active" binding:"omitempty"`
}

// AttributeResponse 属性响应结构
type AttributeResponse struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	DisplayName  string            `json:"display_name"`
	Description  string            `json:"description"`
	Type         AttributeType     `json:"type"`
	Unit         string            `json:"unit"`
	IsRequired   bool              `json:"is_required"`
	IsMultiple   bool              `json:"is_multiple"` // 通过方法计算得出
	DefaultValue string            `json:"default_value"`
	Options      []AttributeOption `json:"options"`
	Validation   ValidationRule    `json:"validation"`
	Sort         int               `json:"sort"`
	IsActive     bool              `json:"is_active"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// CategoryAttributeResponse 分类属性关联响应结构
type CategoryAttributeResponse struct {
	ID          uint              `json:"id"`
	CategoryID  uint              `json:"category_id"`
	AttributeID uint              `json:"attribute_id"`
	IsRequired  bool              `json:"is_required"`
	Sort        int               `json:"sort"`
	Attribute   AttributeResponse `json:"attribute"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// AttributeValueResponse 属性值响应结构
type AttributeValueResponse struct {
	ID          uint               `json:"id"`
	AttributeID uint               `json:"attribute_id"`
	EntityType  string             `json:"entity_type"`
	EntityID    uint               `json:"entity_id"`
	Value       interface{}        `json:"value"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Attribute   *AttributeResponse `json:"attribute,omitempty"`
}

// AttributeListResponse 属性列表响应结构
type AttributeListResponse struct {
	Attributes []AttributeResponse `json:"attributes"`
	Pagination Pagination          `json:"pagination"`
}

// Pagination 分页结构
type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// AttributeQueryRequest 属性查询请求结构
type AttributeQueryRequest struct {
	Name     string        `form:"name"`      // 按名称模糊搜索
	Type     AttributeType `form:"type"`      // 按类型筛选
	IsActive *bool         `form:"is_active"` // 按状态筛选
	Page     int           `form:"page"`
	Limit    int           `form:"limit"`
}

// SetAttributeValueRequest 设置属性值请求结构
type SetAttributeValueRequest struct {
	AttributeID uint        `json:"attribute_id" binding:"required"`
	EntityType  string      `json:"entity_type" binding:"required"`
	EntityID    uint        `json:"entity_id" binding:"required"`
	Value       interface{} `json:"value" binding:"required"`
}

// GetAttributeValuesRequest 获取属性值请求结构
type GetAttributeValuesRequest struct {
	EntityType string `form:"entity_type" binding:"required"`
	EntityID   uint   `form:"entity_id" binding:"required"`
}

// EntityAttributeValuesResponse 实体属性值响应结构
type EntityAttributeValuesResponse struct {
	EntityType string                   `json:"entity_type"`
	EntityID   uint                     `json:"entity_id"`
	Values     []AttributeValueResponse `json:"values"`
}

// 分类属性管理相关结构

// BindAttributeToCategoryRequest 绑定属性到分类请求结构
type BindAttributeToCategoryRequest struct {
	CategoryID  uint `json:"category_id" binding:"required"`
	AttributeID uint `json:"attribute_id" binding:"required"`
	IsRequired  bool `json:"is_required"`
	Sort        int  `json:"sort" binding:"omitempty,min=0"`
}

// UnbindAttributeFromCategoryRequest 从分类解绑属性请求结构
type UnbindAttributeFromCategoryRequest struct {
	CategoryID  uint `json:"category_id" binding:"required"`
	AttributeID uint `json:"attribute_id" binding:"required"`
}

// UpdateCategoryAttributeRequest 更新分类属性关联请求结构
type UpdateCategoryAttributeRequest struct {
	IsRequired *bool `json:"is_required" binding:"omitempty"`
	Sort       *int  `json:"sort" binding:"omitempty,min=0"`
}

// BatchBindAttributesToCategoryRequest 批量绑定属性到分类请求结构
type BatchBindAttributesToCategoryRequest struct {
	CategoryID uint                           `json:"category_id" binding:"required"`
	Attributes []CategoryAttributeBindRequest `json:"attributes" binding:"required"`
}

// CategoryAttributeBindRequest 分类属性绑定请求
type CategoryAttributeBindRequest struct {
	AttributeID uint `json:"attribute_id" binding:"required"`
	IsRequired  bool `json:"is_required"`
	Sort        int  `json:"sort" binding:"omitempty,min=0"`
}

// CategoryAttributesResponse 分类属性列表响应结构
type CategoryAttributesResponse struct {
	CategoryID uint                        `json:"category_id"`
	Attributes []CategoryAttributeResponse `json:"attributes"`
}

// 分类属性继承相关结构

// CategoryAttributeWithInheritanceResponse 带继承信息的分类属性响应结构
type CategoryAttributeWithInheritanceResponse struct {
	ID            uint              `json:"id"`
	CategoryID    uint              `json:"category_id"`
	AttributeID   uint              `json:"attribute_id"`
	IsRequired    bool              `json:"is_required"`
	Sort          int               `json:"sort"`
	IsInherited   bool              `json:"is_inherited"`   // 是否为继承属性
	InheritedFrom *uint             `json:"inherited_from"` // 继承来源分类ID
	Attribute     AttributeResponse `json:"attribute"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// CategoryAttributesWithInheritanceResponse 分类属性列表响应结构（包括继承）
type CategoryAttributesWithInheritanceResponse struct {
	CategoryID uint                                       `json:"category_id"`
	Attributes []CategoryAttributeWithInheritanceResponse `json:"attributes"`
}

// AttributeInheritancePathResponse 属性继承路径响应结构
type AttributeInheritancePathResponse struct {
	CategoryID  uint                                       `json:"category_id"`
	AttributeID uint                                       `json:"attribute_id"`
	Path        []CategoryAttributeWithInheritanceResponse `json:"path"` // 从根分类到当前分类的继承路径
}

// CategoryAttributeBindSummaryResponse 分类属性绑定摘要响应（用于显示继承概览）
type CategoryAttributeBindSummaryResponse struct {
	CategoryID          uint `json:"category_id"`
	TotalAttributes     int  `json:"total_attributes"`     // 总属性数（包括继承）
	OwnAttributes       int  `json:"own_attributes"`       // 自有属性数
	InheritedAttributes int  `json:"inherited_attributes"` // 继承属性数
}
