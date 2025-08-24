package sample

import (
	"erp/internal/modules/sample/handler"
	"erp/internal/modules/sample/repository"
	"erp/internal/modules/sample/service"

	"gorm.io/gorm"
)

// Module 样品模块
type Module struct {
	Handler *handler.Handler
	Service *service.Service
	Repo    *repository.Repository
}

// NewModule 创建样品模块
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
