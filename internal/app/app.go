package app

import (
	"erp/internal/modules/product"
	"erp/internal/modules/source"
	"erp/internal/modules/tags"
	"erp/internal/modules/user"
	"erp/internal/modules/user/handler"
	"erp/internal/modules/user/repository"
	"erp/pkg/oss"
	"log"

	"gorm.io/gorm"
)

// App 应用管理器
type App struct {
	DB      *gorm.DB
	User    *user.Module
	Product *product.Module
	Source  *source.Module
	Tags    *tags.Module
}

// NewApp 创建应用管理器
func NewApp(db *gorm.DB) *App {
	// 初始化OSS客户端
	if err := oss.InitOSS(); err != nil {
		log.Printf("Warning: OSS客户端初始化失败: %v", err)
	} else {
		log.Println("OSS客户端初始化成功")
	}

	// 创建商品模块
	productModule := product.NewModule(db)
	return &App{
		DB:      db,
		User:    user.NewModule(db),
		Product: productModule,
		Source:  source.NewModule(db),
		Tags:    tags.NewModule(db),
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
