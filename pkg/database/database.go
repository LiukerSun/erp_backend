package database

import (
	"fmt"
	"log"

	"erp/config"
	sampleModel "erp/internal/modules/sample/model"
	storeModel "erp/internal/modules/store/model"
	supplierModel "erp/internal/modules/supplier/model"
	"erp/internal/modules/user/model"
	"erp/pkg/password"

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
		&supplierModel.Supplier{},
		&storeModel.Store{},
		&sampleModel.Sample{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// 检查并创建管理员用户（如果不存在）
	var adminUser model.User
	if err := DB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 管理员用户不存在，创建新用户
			hashedPassword, err := password.Hash("admin")
			if err != nil {
				log.Fatal("Failed to hash admin password:", err)
			}

			adminUser = model.User{
				Username: "admin",
				Email:    "admin@example.com",
				Password: hashedPassword,
				Role:     "admin",
			}

			if err := DB.Create(&adminUser).Error; err != nil {
				log.Fatal("Failed to create admin user:", err)
			}
			log.Println("Admin user created successfully")
		} else {
			log.Fatal("Failed to check admin user:", err)
		}
	} else {
		log.Println("Admin user already exists")
	}
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
