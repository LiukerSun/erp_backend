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

		// 分类相关接口
		setupCategoryRoutes(api, app.GetCategoryHandler(), app.GetUserRepository())

		// 产品相关接口
		setupProductRoutes(api, app.GetProductHandler(), app.GetUserRepository())

		// 这里可以添加其他模块的路由
		// setupOrderRoutes(api, app.GetOrderHandler())
		// setupInventoryRoutes(api, app.GetInventoryHandler())
	}
}

// setupUserRoutes 设置用户相关路由
func setupUserRoutes(api *gin.RouterGroup, userHandler interface{}, userRepo interface{}) {
	user := api.Group("/user")
	{
		// 公开接口（无需认证）
		user.POST("/register", userHandler.(interface{ Register(*gin.Context) }).Register)
		user.POST("/login", userHandler.(interface{ Login(*gin.Context) }).Login)

		// 需要认证的接口（包含密码版本验证）
		auth := user.Group("")
		auth.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)))
		{
			auth.GET("/profile", userHandler.(interface{ GetProfile(*gin.Context) }).GetProfile)
			auth.PUT("/profile", userHandler.(interface{ UpdateProfile(*gin.Context) }).UpdateProfile)
			auth.POST("/change_password", userHandler.(interface{ ChangePassword(*gin.Context) }).ChangePassword)
		}

		// 管理员功能路由组（统一管理所有管理员权限相关的接口）
		admin := user.Group("/admin")
		admin.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)), middleware.RoleMiddleware("admin"))
		{
			// 用户列表查询
			admin.GET("/users", userHandler.(interface{ GetUsers(*gin.Context) }).GetUsers)

			// 用户管理操作
			admin.POST("/users", userHandler.(interface{ AdminCreateUser(*gin.Context) }).AdminCreateUser)
			admin.PUT("/users/:id", userHandler.(interface{ AdminUpdateUser(*gin.Context) }).AdminUpdateUser)
			admin.POST("/users/:id/reset_password", userHandler.(interface{ AdminResetUserPassword(*gin.Context) }).AdminResetUserPassword)

			// 删除用户路由，添加防自删除中间件
			admin.DELETE("/users/:id", middleware.PreventSelfDeletionMiddleware(), userHandler.(interface{ AdminDeleteUser(*gin.Context) }).AdminDeleteUser)
		}
	}
}

// setupProductRoutes 设置产品相关路由
func setupProductRoutes(api *gin.RouterGroup, productHandler interface{}, userRepo interface{}) {
	product := api.Group("/product")
	// 产品接口需要认证
	product.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)))
	{
		// 产品 CRUD 操作
		product.POST("", productHandler.(interface{ CreateProduct(*gin.Context) }).CreateProduct)
		product.GET("", productHandler.(interface{ GetProducts(*gin.Context) }).GetProducts)
		product.GET("/:id", productHandler.(interface{ GetProduct(*gin.Context) }).GetProduct)
		product.PUT("/:id", productHandler.(interface{ UpdateProduct(*gin.Context) }).UpdateProduct)
		product.DELETE("/:id", productHandler.(interface{ DeleteProduct(*gin.Context) }).DeleteProduct)

		// 根据分类获取产品
		product.GET("/category/:category_id", productHandler.(interface{ GetProductsByCategory(*gin.Context) }).GetProductsByCategory)
	}
}

// setupCategoryRoutes 设置分类相关路由
func setupCategoryRoutes(api *gin.RouterGroup, categoryHandler interface{}, userRepo interface{}) {
	category := api.Group("/category")
	// 分类接口需要认证
	category.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)))
	{
		// 分类树结构接口
		category.GET("/tree", categoryHandler.(interface{ GetCategoryTree(*gin.Context) }).GetCategoryTree)
		category.GET("/root", categoryHandler.(interface{ GetRootCategories(*gin.Context) }).GetRootCategories)
		category.GET("/:id/children", categoryHandler.(interface{ GetChildrenCategories(*gin.Context) }).GetChildrenCategories)

		// 分类 CRUD 操作
		category.POST("", categoryHandler.(interface{ CreateCategory(*gin.Context) }).CreateCategory)
		category.GET("", categoryHandler.(interface{ GetCategories(*gin.Context) }).GetCategories)
		category.GET("/:id", categoryHandler.(interface{ GetCategory(*gin.Context) }).GetCategory)
		category.PUT("/:id", categoryHandler.(interface{ UpdateCategory(*gin.Context) }).UpdateCategory)
		category.DELETE("/:id", categoryHandler.(interface{ DeleteCategory(*gin.Context) }).DeleteCategory)

		// 分类移动操作
		category.POST("/:id/move", categoryHandler.(interface{ MoveCategory(*gin.Context) }).MoveCategory)
	}
}
