package repository

import (
	"erp/internal/modules/supplier/model"

	"gorm.io/gorm"
)

// Repository 供应商仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建供应商仓库
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建供应商
func (r *Repository) Create(supplier *model.Supplier) error {
	return r.db.Create(supplier).Error
}

// GetByID 根据ID获取供应商
func (r *Repository) GetByID(id uint, includeStores bool) (*model.Supplier, error) {
	var supplier model.Supplier
	query := r.db

	if includeStores {
		query = query.Preload("Stores")
	}

	err := query.First(&supplier, id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

// Update 更新供应商
func (r *Repository) Update(supplier *model.Supplier) error {
	return r.db.Save(supplier).Error
}

// Delete 删除供应商（软删除）
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&model.Supplier{}, id).Error
}

// List 获取供应商列表（支持分页、搜索、排序、筛选）
func (r *Repository) List(page, limit int, search string, includeStores bool, isActive *bool, orderBy string) ([]model.Supplier, int64, error) {
	var suppliers []model.Supplier
	var total int64

	query := r.db.Model(&model.Supplier{})

	// 是否预加载店铺数据
	if includeStores {
		query = query.Preload("Stores")
	}

	// 搜索功能
	if search != "" {
		query = query.Where("name LIKE ? OR remark LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 按活跃状态筛选
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
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

	err = query.Find(&suppliers).Error
	if err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

// CheckNameExists 检查供应商名称是否存在
func (r *Repository) CheckNameExists(name string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Supplier{}).Where("name = ?", name)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}
