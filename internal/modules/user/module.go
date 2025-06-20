package user

import (
	"erp/internal/modules/user/handler"
	"erp/internal/modules/user/repository"
	"erp/internal/modules/user/service"

	"gorm.io/gorm"
)

// Module 用户模块
type Module struct {
	Handler *handler.Handler
	Service *service.Service
	Repo    *repository.Repository
}

// NewModule 创建用户模块
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
