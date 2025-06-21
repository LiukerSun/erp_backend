package routes

import (
	"erp/internal/app"
	"erp/internal/modules/user/repository"
	"erp/pkg/middleware"
	"erp/pkg/oss"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine, app *app.App) {
	// API 路由组
	api := r.Group("/api")
	{
		// 用户相关接口
		setupUserRoutes(api, app.GetUserHandler(), app.GetUserRepository())

		// 商品相关接口
		setupProductRoutes(api, app.Product.GetHandler(), app.GetUserRepository())

		// 货源相关接口
		setupSourceRoutes(api, app.Source.GetHandler(), app.GetUserRepository())

		// 标签相关接口
		setupTagsRoutes(api, app.Tags.GetHandler(), app.GetUserRepository())

		// OSS相关接口
		setupOSSRoutes(api, app.GetUserRepository())
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

// setupProductRoutes 设置商品相关路由
func setupProductRoutes(api *gin.RouterGroup, productHandler interface{}, userRepo interface{}) {
	product := api.Group("/product")
	{
		// 需要认证的接口
		auth := product.Group("")
		auth.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)))
		{
			// 商品管理
			auth.POST("", productHandler.(interface{ Create(*gin.Context) }).Create)
			auth.GET("", productHandler.(interface{ List(*gin.Context) }).List)
			auth.GET("/:id", productHandler.(interface{ Get(*gin.Context) }).Get)
			auth.PUT("/:id", productHandler.(interface{ Update(*gin.Context) }).Update)
			auth.DELETE("/:id", productHandler.(interface{ Delete(*gin.Context) }).Delete)

			// 通过商品编码获取商品
			auth.GET("/code/:code", productHandler.(interface{ GetByCode(*gin.Context) }).GetByCode)
			// 通过SKU获取商品
			auth.GET("/sku/:sku", productHandler.(interface{ GetBySKU(*gin.Context) }).GetBySKU)

			// 图片管理
			auth.PUT("/:id/images/order", productHandler.(interface{ UpdateImageOrder(*gin.Context) }).UpdateImageOrder)
			auth.PUT("/:id/images/main", productHandler.(interface{ SetMainImage(*gin.Context) }).SetMainImage)

			// 颜色管理
			auth.POST("/colors", productHandler.(interface{ CreateColor(*gin.Context) }).CreateColor)
			auth.GET("/colors", productHandler.(interface{ ListColors(*gin.Context) }).ListColors)
			auth.GET("/colors/:id", productHandler.(interface{ GetColor(*gin.Context) }).GetColor)
			auth.PUT("/colors/:id", productHandler.(interface{ UpdateColor(*gin.Context) }).UpdateColor)
			auth.DELETE("/colors/:id", productHandler.(interface{ DeleteColor(*gin.Context) }).DeleteColor)
		}
	}
}

// setupSourceRoutes 设置货源相关路由
func setupSourceRoutes(api *gin.RouterGroup, sourceHandler interface{}, userRepo interface{}) {
	source := api.Group("/source")
	{
		// 需要认证的接口
		auth := source.Group("")
		auth.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)))
		{
			// 货源基本操作
			auth.POST("", sourceHandler.(interface{ Create(*gin.Context) }).Create)
			auth.GET("", sourceHandler.(interface{ List(*gin.Context) }).List)
			auth.GET("/:id", sourceHandler.(interface{ Get(*gin.Context) }).Get)
			auth.PUT("/:id", sourceHandler.(interface{ Update(*gin.Context) }).Update)
			auth.DELETE("/:id", sourceHandler.(interface{ Delete(*gin.Context) }).Delete)

			// 获取启用状态的货源列表
			auth.GET("/active", sourceHandler.(interface{ ListActive(*gin.Context) }).ListActive)
		}
	}
}

// setupTagsRoutes 设置标签相关路由
func setupTagsRoutes(api *gin.RouterGroup, tagsHandler interface{}, userRepo interface{}) {
	tags := api.Group("/tags")
	{
		// 需要认证的接口
		auth := tags.Group("")
		auth.Use(middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)))
		{
			// 标签基本操作
			auth.POST("", tagsHandler.(interface{ CreateTag(*gin.Context) }).CreateTag)
			auth.GET("", tagsHandler.(interface{ GetAllTags(*gin.Context) }).GetAllTags)
			auth.GET("/enabled", tagsHandler.(interface{ GetEnabledTags(*gin.Context) }).GetEnabledTags)
			auth.GET("/:id", tagsHandler.(interface{ GetTagByID(*gin.Context) }).GetTagByID)
			auth.PUT("/:id", tagsHandler.(interface{ UpdateTag(*gin.Context) }).UpdateTag)
			auth.DELETE("/:id", tagsHandler.(interface{ DeleteTag(*gin.Context) }).DeleteTag)

			// 标签与产品关联操作
			auth.GET("/:id/products", tagsHandler.(interface{ GetProductsByTag(*gin.Context) }).GetProductsByTag)
			auth.POST("/:id/products", tagsHandler.(interface{ AddProductToTag(*gin.Context) }).AddProductToTag)
			auth.DELETE("/:id/products", tagsHandler.(interface{ RemoveProductFromTag(*gin.Context) }).RemoveProductFromTag)

			// 获取产品的标签
			auth.GET("/product", tagsHandler.(interface{ GetTagsByProduct(*gin.Context) }).GetTagsByProduct)
		}
	}
}

// setupOSSRoutes 设置OSS相关路由
func setupOSSRoutes(api *gin.RouterGroup, userRepo interface{}) {
	ossGroup := api.Group("/oss")
	{
		// 获取STS临时凭证 (用于前端直传)
		// 这个接口需要认证，确保只有登录用户才能获取上传凭证
		ossGroup.GET("/sts/token", middleware.AuthMiddlewareWithPasswordValidation(userRepo.(*repository.Repository)), oss.GetSTSTokenHandler)
	}
}
