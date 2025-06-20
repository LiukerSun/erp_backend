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

		// 这里可以添加其他模块的路由
		// setupProductRoutes(api, app.GetProductHandler())
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
