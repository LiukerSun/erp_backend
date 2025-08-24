package routes

import (
	"erp/internal/app"
	"erp/internal/modules/user/repository"
	"erp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine, app *app.App) {
	// API 路由组
	api := r.Group("/api")
	{
		// 用户相关接口
		setupUserRoutes(api, app.GetUserHandler(), app.GetUserRepository())
		// 供应商相关接口
		setupSupplierRoutes(api, app.GetSupplierHandler())
		// 店铺相关接口
		setupStoreRoutes(api, app.GetStoreHandler())
		// 样品相关接口
		setupSampleRoutes(api, app.GetSampleHandler())
		// Excel相关接口
		setupExcelRoutes(api, app.GetExcelHandler())
	}
}

// setupSampleRoutes 设置样品相关路由
func setupSampleRoutes(api *gin.RouterGroup, sampleHandler interface{}) {
	sample := api.Group("/samples")
	{
		// 需要认证的接口
		auth := sample.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 获取样品列表
			auth.GET("", sampleHandler.(interface{ ListSamples(*gin.Context) }).ListSamples)
			// 获取样品详情
			auth.GET("/:id", sampleHandler.(interface{ GetSample(*gin.Context) }).GetSample)
			// 创建样品
			auth.POST("", sampleHandler.(interface{ CreateSample(*gin.Context) }).CreateSample)
			// 更新样品
			auth.PUT("/:id", sampleHandler.(interface{ UpdateSample(*gin.Context) }).UpdateSample)
			// 删除样品
			auth.DELETE("/:id", sampleHandler.(interface{ DeleteSample(*gin.Context) }).DeleteSample)
			// 批量更新样品状态
			auth.PATCH("/batch-update", sampleHandler.(interface{ BatchUpdateSamples(*gin.Context) }).BatchUpdateSamples)
		}
	}
}

// setupStoreRoutes 设置店铺相关路由
func setupStoreRoutes(api *gin.RouterGroup, storeHandler interface{}) {
	store := api.Group("/stores")
	{
		// 需要认证的接口
		auth := store.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 获取店铺列表
			auth.GET("", storeHandler.(interface{ ListStores(*gin.Context) }).ListStores)
			// 获取店铺详情
			auth.GET("/:id", storeHandler.(interface{ GetStore(*gin.Context) }).GetStore)
			// 创建店铺
			auth.POST("", storeHandler.(interface{ CreateStore(*gin.Context) }).CreateStore)
			// 更新店铺
			auth.PUT("/:id", storeHandler.(interface{ UpdateStore(*gin.Context) }).UpdateStore)
			// 删除店铺
			auth.DELETE("/:id", storeHandler.(interface{ DeleteStore(*gin.Context) }).DeleteStore)
		}
	}
}

// setupSupplierRoutes 设置供应商相关路由
func setupSupplierRoutes(api *gin.RouterGroup, supplierHandler interface{}) {
	supplier := api.Group("/suppliers")
	{
		// 需要认证的接口
		auth := supplier.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 获取供应商列表
			auth.GET("", supplierHandler.(interface{ ListSuppliers(*gin.Context) }).ListSuppliers)
			// 获取供应商详情
			auth.GET("/:id", supplierHandler.(interface{ GetSupplier(*gin.Context) }).GetSupplier)
			// 创建供应商
			auth.POST("", supplierHandler.(interface{ CreateSupplier(*gin.Context) }).CreateSupplier)
			// 更新供应商
			auth.PUT("/:id", supplierHandler.(interface{ UpdateSupplier(*gin.Context) }).UpdateSupplier)
			// 删除供应商
			auth.DELETE("/:id", supplierHandler.(interface{ DeleteSupplier(*gin.Context) }).DeleteSupplier)
		}
	}
}

// setupExcelRoutes 设置Excel相关路由
func setupExcelRoutes(api *gin.RouterGroup, excelHandler interface{}) {
	excel := api.Group("/excel")
	{
		// 解析Excel文件（无需认证）
		excel.POST("/parse", excelHandler.(interface{ ParseExcel(*gin.Context) }).ParseExcel)
	}
}

// setupUserRoutes 设置用户相关路由
func setupUserRoutes(api *gin.RouterGroup, userHandler interface{}, userRepo interface{}) {
	user := api.Group("/user")
	{
		// 公开接口（无需认证）
		user.POST("/register", userHandler.(interface{ Register(*gin.Context) }).Register)
		user.POST("/login", userHandler.(interface{ Login(*gin.Context) }).Login)
		user.POST("/refresh", userHandler.(interface{ RefreshToken(*gin.Context) }).RefreshToken)

		// 需要认证的接口（包含密码版本验证）
		auth := user.Group("")
		auth.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)))
		{
			auth.GET("/profile", userHandler.(interface{ GetProfile(*gin.Context) }).GetProfile)
			auth.POST("/change_password", userHandler.(interface{ ChangePassword(*gin.Context) }).ChangePassword)
		}

		// 管理员功能路由组（统一管理所有管理员权限相关的接口）
		admin := user.Group("/admin")
		admin.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)), middleware.RoleMiddleware("admin"))
		{
			// 用户列表查询
			admin.GET("/users", userHandler.(interface{ GetUsers(*gin.Context) }).GetUsers)
			// 管理员创建用户
			admin.POST("/users", userHandler.(interface{ AdminCreateUser(*gin.Context) }).AdminCreateUser)
		}
	}
}
