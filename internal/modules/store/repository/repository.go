package repository

import (
	"erp/internal/modules/store/model"

	"gorm.io/gorm"
)

// Repository 店铺仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建店铺仓库
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建店铺
func (r *Repository) Create(store *model.Store) error {
	return r.db.Create(store).Error
}

// GetByID 根据ID获取店铺
func (r *Repository) GetByID(id uint) (*model.Store, error) {
	var store model.Store
	err := r.db.Preload("Supplier").First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

// Update 更新店铺
func (r *Repository) Update(store *model.Store) error {
	return r.db.Save(store).Error
}

// Delete 删除店铺（软删除）
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&model.Store{}, id).Error
}

// List 获取店铺列表（支持分页、搜索、排序、筛选）
func (r *Repository) List(page, limit int, search string, supplierID *uint, isActive *bool, isFeatured *bool, orderBy string) ([]model.Store, int64, error) {
	var stores []model.Store
	var total int64

	query := r.db.Model(&model.Store{}).Preload("Supplier")

	// 搜索功能
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// 按供应商筛选
	if supplierID != nil {
		query = query.Where("supplier_id = ?", *supplierID)
	}

	// 按活跃状态筛选
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// 按精选状态筛选
	if isFeatured != nil {
		query = query.Where("is_featured = ?", *isFeatured)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 排序
	if orderBy == "" {
		orderBy = "created_at DESC"
	}
	query = query.Order(orderBy)

	// 分页查询
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err = query.Find(&stores).Error
	if err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}
