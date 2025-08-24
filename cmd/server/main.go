package main

import (
	"erp/config"
	"erp/internal/app"
	"erp/pkg/database"
	"erp/pkg/middleware"
	"erp/routes"
	"log"

	_ "erp/docs" // 导入 Swagger 文档

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title ERP 后端系统 API
// @version 1.0
// @description 这是一个基于 Go 语言构建的企业资源规划（ERP）后端系统，使用 JWT 进行用户认证，PostgreSQL 作为数据库。API接口按功能模块分组，提供统一的路径结构。
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description 输入 "Bearer " 加上 JWT token

// @tag.name User
// @tag.description 用户管理相关接口，包括注册、登录、资料管理等

// @tag.name Excel
// @tag.description Excel文件上传和解析相关接口

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	database.InitDatabase()

	// 设置Gin模式
	gin.SetMode(config.AppConfig.ServerMode)

	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(middleware.CORSMiddleware())

	// 创建应用管理器
	app := app.NewApp(database.GetDB())

	// 设置路由
	routes.SetupRoutes(r, app)

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "ERP 后端系统运行正常",
		})
	})

	// 启动服务
	addr := ":" + config.AppConfig.ServerPort
	log.Printf("🚀 服务器启动在 http://localhost%s", addr)
	log.Printf("📚 Swagger 文档地址: http://localhost%s/swagger/index.html", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal("❌ 服务器启动失败: ", err)
	}
}
