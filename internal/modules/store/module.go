package store

import (
	"erp/internal/modules/store/handler"
	storeRepo "erp/internal/modules/store/repository"
	"erp/internal/modules/store/service"
	supplierRepo "erp/internal/modules/supplier/repository"

	"gorm.io/gorm"
)

// Module 店铺模块
type Module struct {
	Handler *handler.Handler
	Service *service.Service
	Repo    *storeRepo.Repository
}

// NewModule 创建店铺模块
func NewModule(db *gorm.DB) *Module {
	repo := storeRepo.NewRepository(db)
	supplierRepo := supplierRepo.NewRepository(db)
	svc := service.NewService(repo, supplierRepo)
	h := handler.NewHandler(svc)

	return &Module{
		Handler: h,
		Service: svc,
		Repo:    repo,
	}
}

// GetHandler 获取处理器
func (m *Module) GetHandler() *handler.Handler {
	return m.Handler
}

// GetRepository 获取仓库
func (m *Module) GetRepository() *storeRepo.Repository {
	return m.Repo
}
