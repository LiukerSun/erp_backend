package database

import (
	"fmt"
	"log"

	"erp/config"
	productModel "erp/internal/modules/product/model"
	sourceModel "erp/internal/modules/source/model"
	tagsModel "erp/internal/modules/tags/model"
	"erp/internal/modules/user/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&model.User{},
		&sourceModel.Source{},
		&productModel.Product{},
		&productModel.Color{},
		&productModel.ProductColor{},
		&tagsModel.Tag{},
		&tagsModel.ProductTag{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
