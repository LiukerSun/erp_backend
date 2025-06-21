package attribute

import (
	"erp/internal/modules/attribute/handler"
	"erp/internal/modules/attribute/repository"
	"erp/internal/modules/attribute/service"
	categoryRepo "erp/internal/modules/category/repository"

	"gorm.io/gorm"
)

// Module 属性模块
type Module struct {
	Handler *handler.Handler
	Service *service.Service
	Repo    *repository.Repository
}

// NewModule 创建属性模块
func NewModule(db *gorm.DB, categoryRepo *categoryRepo.Repository) *Module {
	repo := repository.NewRepository(db, categoryRepo)
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

// GetService 获取服务
func (m *Module) GetService() *service.Service {
	return m.Service
}

// GetRepository 获取仓库
func (m *Module) GetRepository() *repository.Repository {
	return m.Repo
}
