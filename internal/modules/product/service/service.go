package service

import (
	"context"
	"erp/internal/modules/product/model"
	"erp/internal/modules/product/repository"
	sourceRepo "erp/internal/modules/source/repository"
	tagsRepo "erp/internal/modules/tags/repository"
	"errors"
	"log"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *model.Product, colorNames []string, tagIDs []uint) error
	UpdateProduct(ctx context.Context, product *model.Product, colorNames []string, tagIDs []uint) error
	DeleteProduct(ctx context.Context, id uint) error
	GetProduct(ctx context.Context, id uint) (*model.Product, error)
	ListProducts(ctx context.Context, page, pageSize int) ([]model.Product, int64, error)
	ListProductsWithFilter(ctx context.Context, filter repository.ProductListFilter, page, pageSize int) ([]model.Product, int64, error)
	CreateColor(ctx context.Context, name, code, hexColor string) (*model.Color, error)
	UpdateColor(ctx context.Context, id uint, name, code, hexColor string) (*model.Color, error)
	DeleteColor(ctx context.Context, id uint) error
	GetColor(ctx context.Context, id uint) (*model.Color, error)
	ListColors(ctx context.Context, orderBy, orderDir string) ([]model.Color, error)
	GetByCode(ctx context.Context, code string, userID uint) (*model.Product, error)
	GetBySKU(ctx context.Context, sku string) (*model.Product, error)
}

type productService struct {
	repo       repository.ProductRepository
	sourceRepo sourceRepo.SourceRepository
	tagsRepo   *tagsRepo.TagsRepository
}

func NewProductService(repo repository.ProductRepository, sourceRepo sourceRepo.SourceRepository, tagsRepo *tagsRepo.TagsRepository) ProductService {
	return &productService{
		repo:       repo,
		sourceRepo: sourceRepo,
		tagsRepo:   tagsRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, product *model.Product, colorNames []string, tagIDs []uint) error {
	// 如果指定了货源ID，获取货源信息并生成商品编码
	if product.SourceID != nil {
		source, err := s.sourceRepo.FindByID(ctx, *product.SourceID)
		if err != nil {
			return errors.New("货源不存在")
		}
		product.Source = &model.Source{
			ID:     source.ID,
			Name:   source.Name,
			Code:   source.Code,
			Status: source.Status,
		}
		// 生成商品编码
		product.GenerateProductCode()

		// 检查商品编码是否已存在
		if product.ProductCode != "" {
			existing, err := s.repo.FindByProductCode(ctx, product.ProductCode)
			if err == nil && existing != nil {
				return errors.New("商品编码已存在")
			}
		}
	}

	// 处理颜色
	colors, err := s.handleColors(ctx, colorNames)
	if err != nil {
		return err
	}
	product.Colors = colors

	// 检查SKU是否已存在
	existing, err := s.repo.FindBySKU(ctx, product.SKU)
	if err == nil && existing != nil {
		return errors.New("商品SKU已存在")
	}

	// 创建产品
	err = s.repo.Create(ctx, product)
	if err != nil {
		return err
	}

	// 处理标签关联
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			err = s.tagsRepo.AddProductToTag(tagID, product.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *productService) UpdateProduct(ctx context.Context, product *model.Product, colorNames []string, tagIDs []uint) error {
	log.Printf("Service: 开始更新商品 ID=%d, 颜色名称=%v", product.ID, colorNames)

	// 检查商品是否存在
	_, err := s.repo.FindByID(ctx, product.ID)
	if err != nil {
		log.Printf("Service: 商品不存在 ID=%d", product.ID)
		return err
	}

	// 如果指定了货源ID，获取货源信息并生成商品编码
	if product.SourceID != nil {
		source, err := s.sourceRepo.FindByID(ctx, *product.SourceID)
		if err != nil {
			return errors.New("货源不存在")
		}
		product.Source = &model.Source{
			ID:     source.ID,
			Name:   source.Name,
			Code:   source.Code,
			Status: source.Status,
		}
		// 生成商品编码
		product.GenerateProductCode()

		// 检查商品编码是否已存在（排除当前商品）
		if product.ProductCode != "" {
			existing, err := s.repo.FindByProductCode(ctx, product.ProductCode)
			if err == nil && existing != nil && existing.ID != product.ID {
				return errors.New("商品编码已存在")
			}
		}
	}

	// 处理颜色
	log.Printf("Service: 处理颜色关联")
	colors, err := s.handleColors(ctx, colorNames)
	if err != nil {
		log.Printf("Service: 处理颜色失败: %v", err)
		return err
	}
	product.Colors = colors
	log.Printf("Service: 颜色处理完成，颜色数量=%d", len(colors))
	for i, color := range colors {
		log.Printf("Service: 颜色[%d] ID=%d, 名称=%s", i, color.ID, color.Name)
	}

	log.Printf("Service: 调用repository更新商品")
	err = s.repo.Update(ctx, product)
	if err != nil {
		return err
	}

	// 处理标签关联
	if tagIDs != nil {
		// 先清除所有标签关联
		existingTags, err := s.tagsRepo.GetTagsByProduct(product.ID)
		if err == nil {
			for _, tag := range existingTags {
				s.tagsRepo.RemoveProductFromTag(tag.ID, product.ID)
			}
		}

		// 添加新的标签关联
		for _, tagID := range tagIDs {
			err = s.tagsRepo.AddProductToTag(tagID, product.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) GetProduct(ctx context.Context, id uint) (*model.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context, page, pageSize int) ([]model.Product, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}

func (s *productService) ListProductsWithFilter(ctx context.Context, filter repository.ProductListFilter, page, pageSize int) ([]model.Product, int64, error) {
	return s.repo.ListWithFilter(ctx, filter, page, pageSize)
}

func (s *productService) CreateColor(ctx context.Context, name, code, hexColor string) (*model.Color, error) {
	// 检查颜色是否已存在
	existing, err := s.repo.FindColorByName(ctx, name)
	if err == nil && existing != nil {
		return existing, nil
	}

	// 如果没有提供代码，自动生成
	if code == "" {
		code = s.generateColorCode(name)
	}

	color := &model.Color{Name: name, Code: code, HexColor: hexColor}
	err = s.repo.CreateColor(ctx, color)
	if err != nil {
		return nil, err
	}
	return color, nil
}

func (s *productService) UpdateColor(ctx context.Context, id uint, name, code, hexColor string) (*model.Color, error) {
	// 检查颜色是否存在
	existing, err := s.repo.FindColorByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新颜色信息
	existing.Name = name
	existing.Code = code
	existing.HexColor = hexColor

	err = s.repo.UpdateColor(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *productService) DeleteColor(ctx context.Context, id uint) error {
	return s.repo.DeleteColor(ctx, id)
}

func (s *productService) GetColor(ctx context.Context, id uint) (*model.Color, error) {
	return s.repo.FindColorByID(ctx, id)
}

func (s *productService) ListColors(ctx context.Context, orderBy, orderDir string) ([]model.Color, error) {
	return s.repo.ListColors(ctx, orderBy, orderDir)
}

// handleColors 处理颜色列表，确保所有颜色都存在
func (s *productService) handleColors(ctx context.Context, colorNames []string) ([]model.Color, error) {
	var colors []model.Color
	for _, name := range colorNames {
		// 根据颜色名称生成代码（转换为大写英文或拼音缩写）
		code := s.generateColorCode(name)
		color, err := s.CreateColor(ctx, name, code, "")
		if err != nil {
			return nil, err
		}
		colors = append(colors, *color)
	}
	return colors, nil
}

// generateColorCode 根据颜色名称生成代码
func (s *productService) generateColorCode(name string) string {
	// 简单的颜色名称到代码的映射
	colorMap := map[string]string{
		"黑色": "BLACK",
		"白色": "WHITE",
		"红色": "RED",
		"蓝色": "BLUE",
		"绿色": "GREEN",
		"黄色": "YELLOW",
		"紫色": "PURPLE",
		"粉色": "PINK",
		"灰色": "GRAY",
		"橙色": "ORANGE",
		"棕色": "BROWN",
		"银色": "SILVER",
		"金色": "GOLD",
		"透明": "TRANSPARENT",
		"彩色": "MULTICOLOR",
	}

	if code, exists := colorMap[name]; exists {
		return code
	}

	// 如果没有预定义的映射，返回原名称的大写版本
	return name
}

// GetByCode 通过SKU获取商品
func (s *productService) GetByCode(ctx context.Context, code string, userID uint) (*model.Product, error) {
	product, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// GetBySKU 通过SKU获取商品
func (s *productService) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	// 获取商品信息
	product, err := s.repo.FindBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}

	return product, nil
}
