package service

import (
	"erp/internal/modules/sample/model"
	"erp/internal/modules/sample/repository"
	"errors"

	"gorm.io/gorm"
)

// Service 样品服务
type Service struct {
	repo *repository.Repository
}

// NewService 创建样品服务
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// CreateSample 创建样品
func (s *Service) CreateSample(req *model.CreateSampleRequest) (*model.SampleResponse, error) {
	// 检查货号是否已存在（在指定供应商下）
	exists, err := s.repo.CheckItemCodeExists(req.ItemCode, req.SupplierID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("该供应商下货号已存在")
	}

	// 检查供应商是否存在
	supplierExists, err := s.repo.CheckSupplierExists(req.SupplierID)
	if err != nil {
		return nil, err
	}
	if !supplierExists {
		return nil, errors.New("供应商不存在")
	}

	// 检查店铺是否存在
	storeExists, err := s.repo.CheckStoreExists(req.StoreID)
	if err != nil {
		return nil, err
	}
	if !storeExists {
		return nil, errors.New("店铺不存在")
	}

	// 检查店铺是否属于指定供应商
	match, err := s.repo.CheckStoreSupplierMatch(req.StoreID, req.SupplierID)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("店铺不属于指定供应商")
	}

	sample := &model.Sample{
		ItemCode:       req.ItemCode,
		SupplierID:     req.SupplierID,
		StoreID:        req.StoreID,
		HasLink:        false,
		IsOffline:      false,
		CanModifyStock: true,
	}

	// 设置可选字段
	if req.HasLink != nil {
		sample.HasLink = *req.HasLink
	}
	if req.IsOffline != nil {
		sample.IsOffline = *req.IsOffline
	}
	if req.CanModifyStock != nil {
		sample.CanModifyStock = *req.CanModifyStock
	}

	err = s.repo.Create(sample)
	if err != nil {
		return nil, err
	}

	// 重新获取样品数据（包含关联数据）
	sample, err = s.repo.GetByID(sample.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(sample), nil
}

// GetSample 获取样品详情
func (s *Service) GetSample(id uint) (*model.SampleResponse, error) {
	sample, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("样品不存在")
		}
		return nil, err
	}

	return s.toResponse(sample), nil
}

// UpdateSample 更新样品
func (s *Service) UpdateSample(id uint, req *model.UpdateSampleRequest) (*model.SampleResponse, error) {
	sample, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("样品不存在")
		}
		return nil, err
	}

	// 检查货号是否已存在（在指定供应商下，排除当前记录）
	exists, err := s.repo.CheckItemCodeExists(req.ItemCode, req.SupplierID, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("该供应商下货号已存在")
	}

	// 检查供应商是否存在
	supplierExists, err := s.repo.CheckSupplierExists(req.SupplierID)
	if err != nil {
		return nil, err
	}
	if !supplierExists {
		return nil, errors.New("供应商不存在")
	}

	// 检查店铺是否存在
	storeExists, err := s.repo.CheckStoreExists(req.StoreID)
	if err != nil {
		return nil, err
	}
	if !storeExists {
		return nil, errors.New("店铺不存在")
	}

	// 检查店铺是否属于指定供应商
	match, err := s.repo.CheckStoreSupplierMatch(req.StoreID, req.SupplierID)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("店铺不属于指定供应商")
	}

	// 更新字段
	sample.ItemCode = req.ItemCode
	sample.SupplierID = req.SupplierID
	sample.StoreID = req.StoreID
	if req.HasLink != nil {
		sample.HasLink = *req.HasLink
	}
	if req.IsOffline != nil {
		sample.IsOffline = *req.IsOffline
	}
	if req.CanModifyStock != nil {
		sample.CanModifyStock = *req.CanModifyStock
	}

	err = s.repo.Update(sample)
	if err != nil {
		return nil, err
	}

	// 重新获取样品数据（包含关联数据）
	sample, err = s.repo.GetByID(sample.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(sample), nil
}

// DeleteSample 删除样品
func (s *Service) DeleteSample(id uint) error {
	// 检查样品是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("样品不存在")
		}
		return err
	}

	return s.repo.Delete(id)
}

// ListSamples 获取样品列表（支持分页、搜索、排序、筛选）
func (s *Service) ListSamples(page, limit int, search string, supplierID, storeID *uint, hasLink, isOffline, canModifyStock *bool, orderBy string) (*model.SampleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	samples, total, err := s.repo.List(page, limit, search, supplierID, storeID, hasLink, isOffline, canModifyStock, orderBy)
	if err != nil {
		return nil, err
	}

	responses := make([]model.SampleResponse, len(samples))
	for i, sample := range samples {
		responses[i] = *s.toResponse(&sample)
	}

	return &model.SampleListResponse{
		Samples: responses,
		Pagination: model.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

// BatchUpdateSamples 批量更新样品状态
func (s *Service) BatchUpdateSamples(req *model.BatchUpdateSamplesRequest) (*model.BatchUpdateSamplesResponse, error) {
	// 验证请求参数
	if len(req.SampleIDs) == 0 {
		return nil, errors.New("样品ID列表不能为空")
	}

	// 检查是否至少有一个要更新的字段
	if req.HasLink == nil && req.IsOffline == nil && req.CanModifyStock == nil {
		return nil, errors.New("至少需要指定一个要更新的字段")
	}

	// 执行批量更新
	successCount, failedCount, err := s.repo.BatchUpdateSamples(req.SampleIDs, req.HasLink, req.IsOffline, req.CanModifyStock)
	if err != nil {
		return nil, err
	}

	totalCount := len(req.SampleIDs)

	// 构造响应消息
	message := "批量更新完成"
	if failedCount > 0 {
		message = "批量更新完成，部分记录更新失败"
	}

	return &model.BatchUpdateSamplesResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		TotalCount:   totalCount,
		Message:      message,
	}, nil
}

// toResponse 转换为响应结构
func (s *Service) toResponse(sample *model.Sample) *model.SampleResponse {
	// 转换供应商信息
	supplierInfo := model.SupplierInfo{
		ID:       sample.Supplier.ID,
		Name:     sample.Supplier.Name,
		Remark:   sample.Supplier.Remark,
		IsActive: sample.Supplier.IsActive,
	}

	// 转换店铺信息
	storeInfo := model.StoreInfo{
		ID:         sample.Store.ID,
		Name:       sample.Store.Name,
		Remark:     sample.Store.Remark,
		IsActive:   sample.Store.IsActive,
		IsFeatured: sample.Store.IsFeatured,
		SupplierID: sample.Store.SupplierID,
	}

	return &model.SampleResponse{
		ID:             sample.ID,
		ItemCode:       sample.ItemCode,
		SupplierID:     sample.SupplierID,
		StoreID:        sample.StoreID,
		HasLink:        sample.HasLink,
		IsOffline:      sample.IsOffline,
		CanModifyStock: sample.CanModifyStock,
		Supplier:       supplierInfo,
		Store:          storeInfo,
		CreatedAt:      sample.CreatedAt,
		UpdatedAt:      sample.UpdatedAt,
	}
}
