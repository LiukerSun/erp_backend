package source

import (
	"erp/internal/modules/source/handler"
	"erp/internal/modules/source/model"
	"erp/internal/modules/source/repository"
	"erp/internal/modules/source/service"

	"gorm.io/gorm"
)

type Module struct {
	db      *gorm.DB
	handler *handler.SourceHandler
}

func NewModule(db *gorm.DB) *Module {
	// 自动迁移数据库表
	db.AutoMigrate(&model.Source{})

	// 创建依赖
	repo := repository.NewSourceRepository(db)
	svc := service.NewSourceService(repo)
	h := handler.NewSourceHandler(svc)

	return &Module{
		db:      db,
		handler: h,
	}
}

// GetHandler 获取货源处理器
func (m *Module) GetHandler() *handler.SourceHandler {
	return m.handler
}
