package supplier

import (
	"erp/internal/modules/supplier/handler"
	"erp/internal/modules/supplier/repository"
	"erp/internal/modules/supplier/service"

	"gorm.io/gorm"
)

// Module 供应商模块
type Module struct {
	Handler *handler.Handler
	Service *service.Service
	Repo    *repository.Repository
}

// NewModule 创建供应商模块
func NewModule(db *gorm.DB) *Module {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
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
