package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Source 货源模型的引用，避免循环导入
type Source struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status int    `json:"status"`
}

// ProductImages 自定义类型，用于处理 JSON 序列化和反序列化
type ProductImages []ProductImage

// Value 实现 driver.Valuer 接口
func (p ProductImages) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "[]", nil
	}
	return json.Marshal(p)
}

// Scan 实现 sql.Scanner 接口
func (p *ProductImages) Scan(value interface{}) error {
	if value == nil {
		*p = ProductImages{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("无法将值转换为 ProductImages")
	}

	if len(bytes) == 0 || string(bytes) == "null" {
		*p = ProductImages{}
		return nil
	}

	// 尝试解析为数组
	err := json.Unmarshal(bytes, p)
	if err != nil {
		// 如果解析失败，检查是否是空对象或单个对象
		var singleImage ProductImage
		if json.Unmarshal(bytes, &singleImage) == nil && singleImage.URL != "" {
			// 如果是单个有效的 ProductImage 对象，转换为数组
			*p = ProductImages{singleImage}
			return nil
		}

		// 检查是否是空对象 {}
		var obj map[string]interface{}
		if json.Unmarshal(bytes, &obj) == nil && len(obj) == 0 {
			*p = ProductImages{}
			return nil
		}

		// 如果都不是，返回原始错误
		return err
	}

	return nil
}

// ProductImage 商品图片信息
// @Description 商品图片信息
type ProductImage struct {
	URL    string `json:"url" example:"https://example.com/image1.jpg"` // 图片URL
	IsMain bool   `json:"is_main" example:"true"`                       // 是否为主图
	Sort   int    `json:"sort" example:"1"`                             // 排序，数字越小越靠前
	Alt    string `json:"alt,omitempty" example:"iPhone 14 正面照"`        // 图片描述
	Title  string `json:"title,omitempty" example:"产品主图"`               // 图片标题
}

// Tag 标签信息（简化版本，避免循环导入）
type Tag struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	IsEnabled   bool   `json:"is_enabled"`
}

// Product 商品模型
// @Description 商品信息
type Product struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	DeletedAt     *time.Time    `json:"deletedAt,omitempty" gorm:"index"`
	Name          string        `json:"name" gorm:"type:varchar(100);not null" example:"iPhone 14"`                   // 商品名称
	SKU           string        `json:"sku" gorm:"type:varchar(50);not null" example:"IPHONE14-128G-BLACK"`           // 货号
	ProductCode   string        `json:"product_code" gorm:"type:varchar(100)" example:"APPLE001-IPHONE14-128G-BLACK"` // 商品编码（店铺编号-货号）
	SourceID      *uint         `json:"source_id" gorm:"index" example:"1"`                                           // 货源ID
	Source        *Source       `json:"source,omitempty" gorm:"foreignKey:SourceID"`                                  // 关联的货源信息
	Price         float64       `json:"price" gorm:"type:decimal(10,2);not null" example:"6999.00"`                   // 售价
	IsDiscounted  bool          `json:"is_discounted" gorm:"default:false" example:"true"`                            // 是否优惠
	DiscountPrice float64       `json:"discount_price" gorm:"type:decimal(10,2)" example:"6799.00"`                   // 优惠价格
	CostPrice     float64       `json:"cost_price" gorm:"type:decimal(10,2);not null" example:"5999.00"`              // 进货价
	Images        ProductImages `json:"images" gorm:"type:json"`                                                      // 商品图片列表
	Colors        []Color       `json:"colors" gorm:"many2many:product_colors;"`                                      // 颜色列表
	Tags          []Tag         `json:"tags" gorm:"many2many:product_tags;"`                                          // 标签列表
	ShippingTime  string        `json:"shipping_time" gorm:"type:varchar(50)" example:"三天"`                           // 发货时间
	IsEnabled     bool          `json:"is_enabled" gorm:"default:true" example:"true"`                                // 是否启用
}

// GenerateProductCode 生成商品编码：店铺编号-货号
func (p *Product) GenerateProductCode() {
	if p.Source != nil && p.Source.Code != "" {
		p.ProductCode = p.Source.Code + "-" + p.SKU
	}
}

// GetMainImage 获取主图
func (p *Product) GetMainImage() *ProductImage {
	for _, img := range p.Images {
		if img.IsMain {
			return &img
		}
	}
	// 如果没有设置主图，返回第一张图片
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

// GetSortedImages 获取按排序的图片列表
func (p *Product) GetSortedImages() []ProductImage {
	images := make([]ProductImage, len(p.Images))
	copy(images, []ProductImage(p.Images))

	// 按sort字段排序，主图优先
	for i := 0; i < len(images)-1; i++ {
		for j := i + 1; j < len(images); j++ {
			// 主图优先
			if images[j].IsMain && !images[i].IsMain {
				images[i], images[j] = images[j], images[i]
			} else if images[i].IsMain == images[j].IsMain {
				// 同样是主图或非主图时，按sort字段排序
				if images[j].Sort < images[i].Sort {
					images[i], images[j] = images[j], images[i]
				}
			}
		}
	}

	return images
}

// SetMainImage 设置主图（取消其他图片的主图状态）
func (p *Product) SetMainImage(url string) {
	for i := range p.Images {
		if p.Images[i].URL == url {
			p.Images[i].IsMain = true
		} else {
			p.Images[i].IsMain = false
		}
	}
}

// Color 颜色模型
// @Description 商品颜色信息
type Color struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`
	Name      string     `json:"name" gorm:"type:varchar(50);uniqueIndex;not null" example:"黑色"` // 颜色名称
	Code      string     `json:"code" gorm:"type:varchar(20);uniqueIndex" example:"BLACK"`       // 颜色代码
	HexColor  string     `json:"hex_color" gorm:"type:varchar(7)" example:"#000000"`             // 十六进制颜色值
	Products  []Product  `json:"products" gorm:"many2many:product_colors;"`                      // 关联的商品
}

// ProductColor 商品和颜色的多对多关联表
type ProductColor struct {
	ProductID uint `gorm:"primaryKey"`
	ColorID   uint `gorm:"primaryKey"`
	CreatedAt time.Time
}
