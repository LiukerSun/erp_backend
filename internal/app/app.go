package app

import (
	"erp/internal/modules/user"
	"erp/internal/modules/user/handler"
	"erp/internal/modules/user/repository"

	"gorm.io/gorm"
)

// App 应用管理器
type App struct {
	DB   *gorm.DB
	User *user.Module
}

// NewApp 创建应用管理器
func NewApp(db *gorm.DB) *App {
	return &App{
		DB:   db,
		User: user.NewModule(db),
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
