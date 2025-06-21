package repository

import (
	"context"
	"erp/internal/modules/product/model"
	"log"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProductListFilter 商品列表筛选条件
type ProductListFilter struct {
	Name         string   `form:"name"`          // 商品名称搜索
	SKU          string   `form:"sku"`           // SKU搜索
	SourceID     *uint    `form:"source_id"`     // 货源ID筛选
	MinPrice     *float64 `form:"min_price"`     // 最低价格
	MaxPrice     *float64 `form:"max_price"`     // 最高价格
	IsDiscounted *bool    `form:"is_discounted"` // 是否优惠
	IsEnabled    *bool    `form:"is_enabled"`    // 是否启用
	ColorNames   []string `form:"colors"`        // 颜色名称列表
	ProductCode  string   `form:"product_code"`  // 商品编码搜索
	ShippingTime string   `form:"shipping_time"` // 发货时间
	OrderBy      string   `form:"order_by"`      // 排序字段: id, name, sku, product_code, price, discount_price, cost_price, is_discounted, is_enabled, shipping_time, created_at, updated_at
	OrderDir     string   `form:"order_dir"`     // 排序方向: asc, desc
}

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.Product, error)
	List(ctx context.Context, page, pageSize int) ([]model.Product, int64, error)
	ListWithFilter(ctx context.Context, filter ProductListFilter, page, pageSize int) ([]model.Product, int64, error)
	FindBySKU(ctx context.Context, sku string) (*model.Product, error)
	FindByProductCode(ctx context.Context, productCode string) (*model.Product, error)
	CreateColor(ctx context.Context, color *model.Color) error
	UpdateColor(ctx context.Context, color *model.Color) error
	DeleteColor(ctx context.Context, id uint) error
	FindColorByID(ctx context.Context, id uint) (*model.Color, error)
	FindColorByName(ctx context.Context, name string) (*model.Color, error)
	ListColors(ctx context.Context, orderBy, orderDir string) ([]model.Color, error)
	GetByCode(code string) (*model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	log.Printf("Repository: 开始更新商品 ID=%d", product.ID)

	// 使用事务来确保数据一致性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 更新商品基本信息（不包含Colors字段，避免GORM自动处理关联）
		updateData := map[string]interface{}{
			"name":           product.Name,
			"sku":            product.SKU,
			"product_code":   product.ProductCode,
			"source_id":      product.SourceID,
			"price":          product.Price,
			"is_discounted":  product.IsDiscounted,
			"discount_price": product.DiscountPrice,
			"cost_price":     product.CostPrice,
			"is_enabled":     product.IsEnabled,
			"images":         product.Images,
			"shipping_time":  product.ShippingTime,
			"updated_at":     product.UpdatedAt,
		}

		log.Printf("Repository: 更新商品基本信息")
		if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Updates(updateData).Error; err != nil {
			log.Printf("Repository: 更新商品基本信息失败: %v", err)
			return err
		}

		// 2. 处理颜色关联关系
		log.Printf("Repository: 删除现有颜色关联")
		// 先删除现有的颜色关联
		if err := tx.Where("product_id = ?", product.ID).Delete(&model.ProductColor{}).Error; err != nil {
			log.Printf("Repository: 删除颜色关联失败: %v", err)
			return err
		}

		// 如果有新的颜色，创建新的颜色关联
		if len(product.Colors) > 0 {
			log.Printf("Repository: 创建新的颜色关联，颜色数量=%d", len(product.Colors))
			var productColors []model.ProductColor
			for _, color := range product.Colors {
				productColors = append(productColors, model.ProductColor{
					ProductID: product.ID,
					ColorID:   color.ID,
				})
				log.Printf("Repository: 创建颜色关联 ProductID=%d, ColorID=%d", product.ID, color.ID)
			}

			if err := tx.Create(&productColors).Error; err != nil {
				log.Printf("Repository: 创建颜色关联失败: %v", err)
				return err
			}
			log.Printf("Repository: 颜色关联创建成功")
		} else {
			log.Printf("Repository: 没有新的颜色关联需要创建")
		}

		log.Printf("Repository: 商品更新完成")
		return nil
	})
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Preload("Source").Preload("Colors").Preload("Tags").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) List(ctx context.Context, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Source").
		Preload("Colors").
		Preload("Tags").
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) ListWithFilter(ctx context.Context, filter ProductListFilter, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&model.Product{})

	// 商品名称模糊搜索
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	// SKU精确搜索
	if filter.SKU != "" {
		query = query.Where("sku = ?", filter.SKU)
	}

	// 商品编码模糊搜索
	if filter.ProductCode != "" {
		query = query.Where("product_code LIKE ?", "%"+filter.ProductCode+"%")
	}

	// 货源筛选
	if filter.SourceID != nil {
		query = query.Where("source_id = ?", *filter.SourceID)
	}

	// 价格范围筛选
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	// 是否优惠筛选
	if filter.IsDiscounted != nil {
		query = query.Where("is_discounted = ?", *filter.IsDiscounted)
	}

	// 是否启用筛选
	if filter.IsEnabled != nil {
		query = query.Where("is_enabled = ?", *filter.IsEnabled)
	}

	// 发货时间筛选
	if filter.ShippingTime != "" {
		query = query.Where("shipping_time LIKE ?", "%"+filter.ShippingTime+"%")
	}

	// 颜色筛选 - 使用子查询
	if len(filter.ColorNames) > 0 {
		colorSubQuery := r.db.Model(&model.Color{}).Select("id").Where("name IN ?", filter.ColorNames)
		query = query.Where("id IN (?)",
			r.db.Model(&model.ProductColor{}).Select("product_id").Where("color_id IN (?)", colorSubQuery))
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := "id"
	orderDir := "desc"

	if filter.OrderBy != "" {
		// 支持所有产品字段排序
		allowedFields := map[string]bool{
			"id":             true,
			"name":           true,
			"sku":            true,
			"product_code":   true,
			"price":          true,
			"discount_price": true,
			"cost_price":     true,
			"is_discounted":  true,
			"is_enabled":     true,
			"shipping_time":  true,
			"created_at":     true,
			"updated_at":     true,
		}

		if allowedFields[filter.OrderBy] {
			orderBy = filter.OrderBy
		}
	}

	if filter.OrderDir != "" {
		switch filter.OrderDir {
		case "asc", "desc":
			orderDir = filter.OrderDir
		}
	}

	// 调试日志
	println("DEBUG: OrderBy =", filter.OrderBy, "OrderDir =", filter.OrderDir)
	println("DEBUG: Final orderBy =", orderBy, "orderDir =", orderDir)
	println("DEBUG: SQL ORDER BY =", orderBy+" "+orderDir)

	// 验证字段是否存在
	if filter.OrderBy != "" {
		allowedFields := map[string]bool{
			"id":             true,
			"name":           true,
			"sku":            true,
			"product_code":   true,
			"price":          true,
			"discount_price": true,
			"cost_price":     true,
			"is_discounted":  true,
			"is_enabled":     true,
			"shipping_time":  true,
			"created_at":     true,
			"updated_at":     true,
		}
		println("DEBUG: Is", filter.OrderBy, "allowed?", allowedFields[filter.OrderBy])
	}

	// 打印完整的SQL查询字符串
	sqlQuery := "SELECT * FROM products ORDER BY " + orderBy + " " + orderDir + " LIMIT " + strconv.Itoa(pageSize) + " OFFSET " + strconv.Itoa((page-1)*pageSize)
	println("DEBUG: Expected SQL:", sqlQuery)

	// 执行查询
	// 临时启用调试模式查看SQL
	debugDB := query.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)})

	err = debugDB.
		Order(orderBy + " " + orderDir).
		Preload("Source").
		Preload("Colors").
		Preload("Tags").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&products).Error

	// 调试日志 - 查看生成的SQL
	if err == nil {
		println("DEBUG: Query executed successfully")
		println("DEBUG: Found", len(products), "products")
	} else {
		println("DEBUG: Query error:", err.Error())
	}

	return products, total, err
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByProductCode(ctx context.Context, productCode string) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Where("product_code = ?", productCode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) CreateColor(ctx context.Context, color *model.Color) error {
	return r.db.WithContext(ctx).Create(color).Error
}

func (r *productRepository) UpdateColor(ctx context.Context, color *model.Color) error {
	return r.db.WithContext(ctx).Save(color).Error
}

func (r *productRepository) DeleteColor(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Color{}, id).Error
}

func (r *productRepository) FindColorByID(ctx context.Context, id uint) (*model.Color, error) {
	var color model.Color
	err := r.db.WithContext(ctx).First(&color, id).Error
	if err != nil {
		return nil, err
	}
	return &color, nil
}

func (r *productRepository) FindColorByName(ctx context.Context, name string) (*model.Color, error) {
	var color model.Color
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&color).Error
	if err != nil {
		return nil, err
	}
	return &color, nil
}

func (r *productRepository) ListColors(ctx context.Context, orderBy, orderDir string) ([]model.Color, error) {
	var colors []model.Color

	// 验证排序字段
	allowedFields := map[string]bool{
		"id":         true,
		"name":       true,
		"code":       true,
		"hex_color":  true,
		"created_at": true,
		"updated_at": true,
	}

	// 默认排序
	sortField := "id"
	sortDirection := "ASC"

	if orderBy != "" && allowedFields[orderBy] {
		sortField = orderBy
	}

	if orderDir != "" {
		switch orderDir {
		case "asc", "ASC":
			sortDirection = "ASC"
		case "desc", "DESC":
			sortDirection = "DESC"
		}
	}

	err := r.db.WithContext(ctx).Order(sortField + " " + sortDirection).Find(&colors).Error
	return colors, err
}

// GetByCode 通过SKU获取商品
func (r *productRepository) GetByCode(code string) (*model.Product, error) {
	var product model.Product
	if err := r.db.Preload("Source").Preload("Colors").Preload("Tags").Where("sku = ?", code).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
