package tags

import (
	"erp/internal/modules/tags/handler"
	"erp/internal/modules/tags/model"
	"erp/internal/modules/tags/repository"
	"erp/internal/modules/tags/service"

	"gorm.io/gorm"
)

type Module struct {
	db      *gorm.DB
	handler *handler.TagsHandler
}

func NewModule(db *gorm.DB) *Module {
	// 自动迁移数据库表
	db.AutoMigrate(&model.Tag{}, &model.ProductTag{})

	// 创建依赖
	repo := repository.NewTagsRepository(db)
	svc := service.NewTagsService(repo)
	handler := handler.NewTagsHandler(svc)

	return &Module{
		db:      db,
		handler: handler,
	}
}

// GetHandler 获取标签处理器
func (m *Module) GetHandler() *handler.TagsHandler {
	return m.handler
}
