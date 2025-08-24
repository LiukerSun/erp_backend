package service

import (
	"erp/internal/modules/supplier/model"
	"erp/internal/modules/supplier/repository"
	"errors"

	"gorm.io/gorm"
)

// Service 供应商服务
type Service struct {
	repo *repository.Repository
}

// NewService 创建供应商服务
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// CreateSupplier 创建供应商
func (s *Service) CreateSupplier(req *model.CreateSupplierRequest) (*model.SupplierResponse, error) {
	// 检查供应商名称是否已存在
	exists, err := s.repo.CheckNameExists(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("供应商名称已存在")
	}

	supplier := &model.Supplier{
		Name:   req.Name,
		Remark: req.Remark,
	}

	err = s.repo.Create(supplier)
	if err != nil {
		return nil, err
	}

	return s.toResponse(supplier), nil
}

// GetSupplier 获取供应商详情
func (s *Service) GetSupplier(id uint, includeStores bool) (*model.SupplierResponse, error) {
	supplier, err := s.repo.GetByID(id, includeStores)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("供应商不存在")
		}
		return nil, err
	}

	return s.toResponse(supplier), nil
}

// UpdateSupplier 更新供应商
func (s *Service) UpdateSupplier(id uint, req *model.UpdateSupplierRequest) (*model.SupplierResponse, error) {
	supplier, err := s.repo.GetByID(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("供应商不存在")
		}
		return nil, err
	}

	// 检查供应商名称是否已存在（排除当前记录）
	exists, err := s.repo.CheckNameExists(req.Name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("供应商名称已存在")
	}

	// 更新字段
	supplier.Name = req.Name
	supplier.Remark = req.Remark
	if req.IsActive != nil {
		supplier.IsActive = *req.IsActive
	}

	err = s.repo.Update(supplier)
	if err != nil {
		return nil, err
	}

	return s.toResponse(supplier), nil
}

// DeleteSupplier 删除供应商
func (s *Service) DeleteSupplier(id uint) error {
	// 检查供应商是否存在
	_, err := s.repo.GetByID(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("供应商不存在")
		}
		return err
	}

	return s.repo.Delete(id)
}

// ListSuppliers 获取供应商列表（支持分页、搜索、排序、筛选）
func (s *Service) ListSuppliers(page, limit int, search string, includeStores bool, isActive *bool, orderBy string) (*model.SupplierListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	suppliers, total, err := s.repo.List(page, limit, search, includeStores, isActive, orderBy)
	if err != nil {
		return nil, err
	}

	responses := make([]model.SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		responses[i] = *s.toResponse(&supplier)
	}

	return &model.SupplierListResponse{
		Suppliers: responses,
		Pagination: model.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

// toResponse 转换为响应结构
func (s *Service) toResponse(supplier *model.Supplier) *model.SupplierResponse {
	// 转换店铺信息
	storeInfos := make([]model.StoreInfo, len(supplier.Stores))
	for i, store := range supplier.Stores {
		storeInfos[i] = model.StoreInfo{
			ID:         store.ID,
			Name:       store.Name,
			Remark:     store.Remark,
			IsActive:   store.IsActive,
			SupplierID: store.SupplierID,
			CreatedAt:  store.CreatedAt,
			UpdatedAt:  store.UpdatedAt,
		}
	}

	return &model.SupplierResponse{
		ID:        supplier.ID,
		Name:      supplier.Name,
		Remark:    supplier.Remark,
		IsActive:  supplier.IsActive,
		Stores:    storeInfos,
		CreatedAt: supplier.CreatedAt,
		UpdatedAt: supplier.UpdatedAt,
	}
}
