package app

import (
	"erp/internal/modules/category"
	categoryHandler "erp/internal/modules/category/handler"
	"erp/internal/modules/product"
	productHandler "erp/internal/modules/product/handler"
	"erp/internal/modules/user"
	"erp/internal/modules/user/handler"
	"erp/internal/modules/user/repository"

	"gorm.io/gorm"
)

// App 应用管理器
type App struct {
	DB       *gorm.DB
	User     *user.Module
	Category *category.Module
	Product  *product.Module
}

// NewApp 创建应用管理器
func NewApp(db *gorm.DB) *App {
	// 先创建分类模块
	categoryModule := category.NewModule(db)

	// 创建产品模块时传入分类仓库依赖
	productModule := product.NewModule(db, categoryModule.GetRepository())

	return &App{
		DB:       db,
		User:     user.NewModule(db),
		Category: categoryModule,
		Product:  productModule,
	}
}

// GetUserHandler 获取用户处理器
func (a *App) GetUserHandler() *handler.Handler {
	return a.User.GetHandler()
}

// GetUserRepository 获取用户仓库
func (a *App) GetUserRepository() *repository.Repository {
	return a.User.GetRepository()
}

// GetCategoryHandler 获取分类处理器
func (a *App) GetCategoryHandler() *categoryHandler.Handler {
	return a.Category.GetHandler()
}

// GetProductHandler 获取产品处理器
func (a *App) GetProductHandler() *productHandler.Handler {
	return a.Product.GetHandler()
}
