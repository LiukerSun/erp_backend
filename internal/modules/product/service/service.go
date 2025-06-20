package service

import (
	"context"
	"errors"

	categoryRepo "erp/internal/modules/category/repository"
	"erp/internal/modules/product/model"
	"erp/internal/modules/product/repository"

	"gorm.io/gorm"
)

// Service 产品服务
type Service struct {
	repo         *repository.Repository
	categoryRepo *categoryRepo.Repository
}

// NewService 创建产品服务
func NewService(repo *repository.Repository, categoryRepo *categoryRepo.Repository) *Service {
	return &Service{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

// CreateProduct 创建产品
func (s *Service) CreateProduct(ctx context.Context, req model.CreateProductRequest) (*model.ProductResponse, error) {
	// 创建产品
	product := &model.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, errors.New("产品创建失败")
	}

	// 返回产品信息
	return &model.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}, nil
}

// GetProduct 获取产品详情
func (s *Service) GetProduct(ctx context.Context, id uint) (*model.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, errors.New("获取产品信息失败")
	}

	return &model.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}, nil
}

// UpdateProduct 更新产品
func (s *Service) UpdateProduct(ctx context.Context, id uint, req model.UpdateProductRequest) (*model.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, errors.New("获取产品信息失败")
	}

	// 更新产品信息
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, errors.New("产品更新失败")
	}

	return &model.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}, nil
}

// DeleteProduct 删除产品
func (s *Service) DeleteProduct(ctx context.Context, id uint) error {
	// 检查产品是否存在
	if !s.repo.ExistsByID(ctx, id) {
		return errors.New("产品不存在")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.New("产品删除失败")
	}

	return nil
}

// SearchProducts 搜索产品（支持筛选和分页）
func (s *Service) SearchProducts(ctx context.Context, req model.ProductQueryRequest) (*model.ProductListResponse, error) {
	// 设置默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	// 使用查询接口，如果没有筛选条件，会返回所有产品
	products, total, err := s.repo.FindWithQuery(ctx, req, offset, req.Limit)
	if err != nil {
		return nil, errors.New("获取产品列表失败")
	}

	// 转换为响应格式
	var productResponses []model.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, model.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			CategoryID: product.CategoryID,
			CreatedAt:  product.CreatedAt,
			UpdatedAt:  product.UpdatedAt,
		})
	}

	return &model.ProductListResponse{
		Products: productResponses,
		Pagination: model.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
		},
	}, nil
}

// GetProductsByCategory 根据分类获取产品
func (s *Service) GetProductsByCategory(ctx context.Context, categoryID uint) ([]model.ProductResponse, error) {
	products, err := s.repo.FindByCategory(ctx, categoryID)
	if err != nil {
		return nil, errors.New("获取分类产品失败")
	}

	// 转换为响应格式
	var productResponses []model.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, model.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			CategoryID: product.CategoryID,
			CreatedAt:  product.CreatedAt,
			UpdatedAt:  product.UpdatedAt,
		})
	}

	return productResponses, nil
}
