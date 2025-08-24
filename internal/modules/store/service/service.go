package service

import (
	"erp/internal/modules/store/model"
	storeRepo "erp/internal/modules/store/repository"
	supplierRepo "erp/internal/modules/supplier/repository"
	"errors"

	"gorm.io/gorm"
)

// Service 店铺服务
type Service struct {
	repo         *storeRepo.Repository
	supplierRepo *supplierRepo.Repository
}

// NewService 创建店铺服务
func NewService(repo *storeRepo.Repository, supplierRepo *supplierRepo.Repository) *Service {
	return &Service{
		repo:         repo,
		supplierRepo: supplierRepo,
	}
}

// CreateStore 创建店铺
func (s *Service) CreateStore(req *model.CreateStoreRequest) (*model.StoreResponse, error) {
	// 验证供应商是否存在
	_, err := s.supplierRepo.GetByID(req.SupplierID, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("供应商不存在")
		}
		return nil, err
	}

	store := &model.Store{
		Name:                  req.Name,
		Remark:                req.Remark,
		SupplierID:            req.SupplierID,
		DefaultCommissionRate: req.DefaultCommissionRate,
	}

	// 处理 IsFeatured 属性
	if req.IsFeatured != nil {
		store.IsFeatured = *req.IsFeatured
	}

	err = s.repo.Create(store)
	if err != nil {
		return nil, err
	}

	return s.toResponse(store), nil
}

// GetStore 获取店铺详情
func (s *Service) GetStore(id uint) (*model.StoreResponse, error) {
	store, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("店铺不存在")
		}
		return nil, err
	}

	return s.toResponse(store), nil
}

// UpdateStore 更新店铺
func (s *Service) UpdateStore(id uint, req *model.UpdateStoreRequest) (*model.StoreResponse, error) {
	store, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("店铺不存在")
		}
		return nil, err
	}

	// 验证供应商是否存在
	_, err = s.supplierRepo.GetByID(req.SupplierID, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("供应商不存在")
		}
		return nil, err
	}

	// 更新字段
	store.Name = req.Name
	store.Remark = req.Remark
	store.SupplierID = req.SupplierID
	store.DefaultCommissionRate = req.DefaultCommissionRate
	if req.IsActive != nil {
		store.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		store.IsFeatured = *req.IsFeatured
	}

	if req.DefaultCommissionRate != 0 {
		store.DefaultCommissionRate = req.DefaultCommissionRate
	}

	err = s.repo.Update(store)
	if err != nil {
		return nil, err
	}

	return s.toResponse(store), nil
}

// DeleteStore 删除店铺
func (s *Service) DeleteStore(id uint) error {
	// 检查店铺是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("店铺不存在")
		}
		return err
	}

	return s.repo.Delete(id)
}

// ListStores 获取店铺列表（支持分页、搜索、排序、筛选）
func (s *Service) ListStores(page, limit int, search string, supplierID *uint, isActive *bool, isFeatured *bool, orderBy string) (*model.StoreListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	stores, total, err := s.repo.List(page, limit, search, supplierID, isActive, isFeatured, orderBy)
	if err != nil {
		return nil, err
	}

	responses := make([]model.StoreResponse, len(stores))
	for i, store := range stores {
		responses[i] = *s.toResponse(&store)
	}

	return &model.StoreListResponse{
		Stores: responses,
		Pagination: model.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

// toResponse 转换为响应结构
func (s *Service) toResponse(store *model.Store) *model.StoreResponse {
	return &model.StoreResponse{
		ID:                    store.ID,
		Name:                  store.Name,
		Remark:                store.Remark,
		IsActive:              store.IsActive,
		IsFeatured:            store.IsFeatured,
		DefaultCommissionRate: store.DefaultCommissionRate,
		SupplierID:            store.SupplierID,
		Supplier:              store.Supplier,
		CreatedAt:             store.CreatedAt,
		UpdatedAt:             store.UpdatedAt,
	}
}
