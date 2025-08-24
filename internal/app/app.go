package app

import (
	"erp/config"
	"erp/internal/modules/excel"
	sampleModule "erp/internal/modules/sample"
	sampleHandler "erp/internal/modules/sample/handler"
	sampleRepo "erp/internal/modules/sample/repository"
	storeModule "erp/internal/modules/store"
	storeHandler "erp/internal/modules/store/handler"
	storeRepo "erp/internal/modules/store/repository"
	supplierModule "erp/internal/modules/supplier"
	supplierHandler "erp/internal/modules/supplier/handler"
	supplierRepo "erp/internal/modules/supplier/repository"
	userModule "erp/internal/modules/user"
	userHandler "erp/internal/modules/user/handler"
	userRepo "erp/internal/modules/user/repository"

	"gorm.io/gorm"
)

// App 应用管理器
type App struct {
	DB       *gorm.DB
	User     *userModule.Module
	Supplier *supplierModule.Module
	Store    *storeModule.Module
	Sample   *sampleModule.Module
	Excel    *excel.Module
}

// NewApp 创建应用管理器
func NewApp(db *gorm.DB) *App {
	// 创建Excel模块
	excelModule, err := excel.NewModule(config.AppConfig.GrpcServerAddr)
	if err != nil {
		panic("创建Excel模块失败: " + err.Error())
	}

	// 创建依赖模块
	supplierModule := supplierModule.NewModule(db)

	return &App{
		DB:       db,
		User:     userModule.NewModule(db),
		Supplier: supplierModule,
		Store:    storeModule.NewModule(db),
		Sample:   sampleModule.NewModule(db),
		Excel:    excelModule,
	}
}

// GetUserHandler 获取用户处理器
func (a *App) GetUserHandler() *userHandler.Handler {
	return a.User.GetHandler()
}

// GetUserRepository 获取用户仓库
func (a *App) GetUserRepository() *userRepo.Repository {
	return a.User.GetRepository()
}

// GetSupplierHandler 获取供应商处理器
func (a *App) GetSupplierHandler() *supplierHandler.Handler {
	return a.Supplier.GetHandler()
}

// GetSupplierRepository 获取供应商仓库
func (a *App) GetSupplierRepository() *supplierRepo.Repository {
	return a.Supplier.GetRepository()
}

// GetStoreHandler 获取店铺处理器
func (a *App) GetStoreHandler() *storeHandler.Handler {
	return a.Store.GetHandler()
}

// GetStoreRepository 获取店铺仓库
func (a *App) GetStoreRepository() *storeRepo.Repository {
	return a.Store.GetRepository()
}

// GetSampleHandler 获取样品处理器
func (a *App) GetSampleHandler() *sampleHandler.Handler {
	return a.Sample.GetHandler()
}

// GetSampleRepository 获取样品仓库
func (a *App) GetSampleRepository() *sampleRepo.Repository {
	return a.Sample.GetRepository()
}

// GetExcelHandler 获取Excel处理器
func (a *App) GetExcelHandler() interface{} {
	return a.Excel.Handler
}

// Close 关闭应用资源
func (a *App) Close() error {
	if a.Excel != nil {
		return a.Excel.Close()
	}
	return nil
}
