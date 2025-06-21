package product

import (
	attributeService "erp/internal/modules/attribute/service"
	categoryRepo "erp/internal/modules/category/repository"
	"erp/internal/modules/product/handler"
	"erp/internal/modules/product/repository"
	"erp/internal/modules/product/service"

	"gorm.io/gorm"
)

// Module 产品模块
type Module struct {
	Handler *handler.Handler
	Service *service.Service
	Repo    *repository.Repository
}

// NewModule 创建产品模块
func NewModule(db *gorm.DB, categoryRepo *categoryRepo.Repository, attributeService *attributeService.Service) *Module {
	repo := repository.NewRepository(db, categoryRepo)
	svc := service.NewService(repo, categoryRepo, attributeService)
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
func (m *Module) GetRepository() *repository.Repository {
	return m.Repo
}
