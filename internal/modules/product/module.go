package product

import (
	"erp/internal/modules/product/handler"
	"erp/internal/modules/product/model"
	"erp/internal/modules/product/repository"
	"erp/internal/modules/product/service"
	sourceRepo "erp/internal/modules/source/repository"
	tagsModel "erp/internal/modules/tags/model"
	tagsRepo "erp/internal/modules/tags/repository"

	"gorm.io/gorm"
)

type Module struct {
	db      *gorm.DB
	handler *handler.ProductHandler
}

func NewModule(db *gorm.DB) *Module {
	// 自动迁移数据库表
	db.AutoMigrate(&model.Product{}, &model.Color{}, &model.ProductColor{}, &tagsModel.Tag{}, &tagsModel.ProductTag{})

	// 创建依赖
	productRepo := repository.NewProductRepository(db)
	sourceRepository := sourceRepo.NewSourceRepository(db)
	tagsRepository := tagsRepo.NewTagsRepository(db)
	svc := service.NewProductService(productRepo, sourceRepository, tagsRepository)
	h := handler.NewProductHandler(svc)

	return &Module{
		db:      db,
		handler: h,
	}
}

// GetHandler 获取商品处理器
func (m *Module) GetHandler() *handler.ProductHandler {
	return m.handler
}
