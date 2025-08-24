package repository

import (
	"erp/internal/modules/sample/model"
	storeModel "erp/internal/modules/store/model"
	supplierModel "erp/internal/modules/supplier/model"

	"gorm.io/gorm"
)

// Repository 样品仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建样品仓库
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建样品
func (r *Repository) Create(sample *model.Sample) error {
	return r.db.Create(sample).Error
}

// GetByID 根据ID获取样品
func (r *Repository) GetByID(id uint) (*model.Sample, error) {
	var sample model.Sample
	err := r.db.Preload("Supplier").Preload("Store").First(&sample, id).Error
	if err != nil {
		return nil, err
	}
	return &sample, nil
}

// Update 更新样品
func (r *Repository) Update(sample *model.Sample) error {
	return r.db.Save(sample).Error
}

// Delete 删除样品（软删除）
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&model.Sample{}, id).Error
}

// List 获取样品列表（支持分页、搜索、排序、筛选）
func (r *Repository) List(page, limit int, search string, supplierID, storeID *uint, hasLink, isOffline, canModifyStock *bool, orderBy string) ([]model.Sample, int64, error) {
	var samples []model.Sample
	var total int64

	query := r.db.Model(&model.Sample{}).Preload("Supplier").Preload("Store")

	// 搜索功能（搜索货号）
	if search != "" {
		query = query.Where("item_code LIKE ?", "%"+search+"%")
	}

	// 按供应商筛选
	if supplierID != nil {
		query = query.Where("supplier_id = ?", *supplierID)
	}

	// 按店铺筛选
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	// 按是否制作链接筛选
	if hasLink != nil {
		query = query.Where("has_link = ?", *hasLink)
	}

	// 按是否下架筛选
	if isOffline != nil {
		query = query.Where("is_offline = ?", *isOffline)
	}

	// 按是否可修改库存筛选
	if canModifyStock != nil {
		query = query.Where("can_modify_stock = ?", *canModifyStock)
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

	err = query.Find(&samples).Error
	if err != nil {
		return nil, 0, err
	}

	return samples, total, nil
}

// CheckItemCodeExists 检查货号是否存在（在指定供应商下）
func (r *Repository) CheckItemCodeExists(itemCode string, supplierID uint, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Sample{}).Where("item_code = ? AND supplier_id = ?", itemCode, supplierID)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// CheckSupplierExists 检查供应商是否存在
func (r *Repository) CheckSupplierExists(supplierID uint) (bool, error) {
	var count int64
	err := r.db.Model(&supplierModel.Supplier{}).Where("id = ?", supplierID).Count(&count).Error
	return count > 0, err
}

// CheckStoreExists 检查店铺是否存在
func (r *Repository) CheckStoreExists(storeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&storeModel.Store{}).Where("id = ?", storeID).Count(&count).Error
	return count > 0, err
}

// CheckStoreSupplierMatch 检查店铺是否属于指定供应商
func (r *Repository) CheckStoreSupplierMatch(storeID, supplierID uint) (bool, error) {
	var count int64
	err := r.db.Model(&storeModel.Store{}).Where("id = ? AND supplier_id = ?", storeID, supplierID).Count(&count).Error
	return count > 0, err
}

// BatchUpdateSamples 批量更新样品状态
func (r *Repository) BatchUpdateSamples(sampleIDs []uint, hasLink, isOffline, canModifyStock *bool) (int, int, error) {
	// 查询要更新的样品
	var samples []model.Sample
	err := r.db.Where("id IN ?", sampleIDs).Find(&samples).Error
	if err != nil {
		return 0, 0, err
	}

	foundCount := len(samples)
	successCount := 0
	failedCount := 0

	// 批量更新
	updates := make(map[string]interface{})
	if hasLink != nil {
		updates["has_link"] = *hasLink
	}
	if isOffline != nil {
		updates["is_offline"] = *isOffline
	}
	if canModifyStock != nil {
		updates["can_modify_stock"] = *canModifyStock
	}

	// 如果有要更新的字段
	if len(updates) > 0 {
		result := r.db.Model(&model.Sample{}).Where("id IN ?", sampleIDs).Updates(updates)
		if result.Error != nil {
			return 0, len(sampleIDs), result.Error
		}
		successCount = int(result.RowsAffected)
	}

	// 计算失败数量
	totalRequested := len(sampleIDs)
	notFoundCount := totalRequested - foundCount
	failedCount = notFoundCount + (foundCount - successCount)

	return successCount, failedCount, nil
}

// GetSamplesByIDs 根据ID列表获取样品
func (r *Repository) GetSamplesByIDs(ids []uint) ([]model.Sample, error) {
	var samples []model.Sample
	err := r.db.Where("id IN ?", ids).Find(&samples).Error
	return samples, err
}
