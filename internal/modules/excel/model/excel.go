package model

import (
	"time"
)

// ExcelUploadRequest Excel上传请求
type ExcelUploadRequest struct {
	SheetName string `json:"sheet_name" form:"sheet_name"` // Sheet名称（可选）
}

// ProductInfo 商品信息
type ProductInfo struct {
	ProductID       string `json:"product_id"`        // 商品ID
	ProductName     string `json:"product_name"`      // 商品名称
	CategoryLevel1  string `json:"category_level1"`   // 一级类目
	CategoryLevel2  string `json:"category_level2"`   // 二级类目
	CategoryLevel3  string `json:"category_level3"`   // 三级类目
	CategoryLevel4  string `json:"category_level4"`   // 四级类目
	ProductType     string `json:"product_type"`      // 商品类型
	ProductGroup    string `json:"product_group"`     // 商品分组
	MerchantCode    string `json:"merchant_code"`     // 商家编码
	MerchantSkuCode string `json:"merchant_sku_code"` // 商家SKU编码
	SpecID          string `json:"spec_id"`           // 规格ID（SKUID）
	ProductSpec     string `json:"product_spec"`      // 商品规格
	Color           string `json:"color"`             // 颜色分类
	Size            string `json:"size"`              // 尺码大小
	ShippingTime    string `json:"shipping_time"`     // 发货时效（从ProductSpec中提取）
	DeliveryTime    string `json:"delivery_time"`     // 商品发货时间（原始字段）
	Price           string `json:"price"`             // 商品价格
	InStock         string `json:"in_stock"`          // 现货可售
	PreSaleStock    string `json:"pre_sale_stock"`    // 预售库存
	TieredStock     string `json:"tiered_stock"`      // 阶梯库存
	SalesVolume     string `json:"sales_volume"`      // 销量
	CommissionRate  string `json:"commission_rate"`   // 佣金比例
	AuditStatus     string `json:"audit_status"`      // 商品审核状态
	ProductLink     string `json:"product_link"`      // 商品链接
	ProductCode     string `json:"product_code"`      // 货号
}

// ExcelParseResponse Excel解析响应
type ExcelParseResponse struct {
	Success  bool          `json:"success"`   // 是否成功
	Message  string        `json:"message"`   // 消息
	Products []ProductInfo `json:"products"`  // 商品列表
	Total    int           `json:"total"`     // 总数
	UploadAt time.Time     `json:"upload_at"` // 上传时间
}
