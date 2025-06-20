package repository

import (
	"context"
	"fmt"

	categoryRepo "erp/internal/modules/category/repository"
	"erp/internal/modules/product/model"

	"gorm.io/gorm"
)

// Repository 产品仓库
type Repository struct {
	db           *gorm.DB
	categoryRepo *categoryRepo.Repository
}

// NewRepository 创建产品仓库
func NewRepository(db *gorm.DB, categoryRepo *categoryRepo.Repository) *Repository {
	return &Repository{
		db:           db,
		categoryRepo: categoryRepo,
	}
}

// Create 创建产品
func (r *Repository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// FindByID 根据ID查找产品
func (r *Repository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update 更新产品
func (r *Repository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete 软删除产品
func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}

// FindWithPagination 分页查找产品
func (r *Repository) FindWithPagination(ctx context.Context, offset, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取产品列表
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// FindWithQuery 根据查询条件查找产品
func (r *Repository) FindWithQuery(ctx context.Context, req model.ProductQueryRequest, offset, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{})

	// 按名称模糊搜索
	if req.Name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Name))
	}

	// 按分类筛选（包含子分类）
	if req.CategoryID > 0 {
		// 获取当前分类及其所有子分类的ID
		categoryIDs := []uint{req.CategoryID}

		// 获取所有子分类
		descendants, err := r.categoryRepo.GetDescendants(req.CategoryID)
		if err != nil {
			// 如果获取子分类失败，仍然使用原分类进行查询
			// 这样可以确保即使分类服务有问题，基本查询仍然可用
		} else {
			// 添加所有子分类ID
			for _, desc := range descendants {
				categoryIDs = append(categoryIDs, desc.ID)
			}
		}

		// 使用IN查询包含所有相关分类的产品
		query = query.Where("category_id IN ?", categoryIDs)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取产品列表
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// FindByCategory 根据分类查找产品
func (r *Repository) FindByCategory(ctx context.Context, categoryID uint) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ExistsByID 检查产品是否存在
func (r *Repository) ExistsByID(ctx context.Context, id uint) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Count(&count)
	return count > 0
}
