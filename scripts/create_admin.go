package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"erp/config"
	"erp/internal/modules/user/model"
	"erp/internal/modules/user/repository"
	"erp/pkg/database"
	"erp/pkg/password"
)

func main() {
	var (
		username = flag.String("username", "admin", "管理员用户名")
		email    = flag.String("email", "admin@example.com", "管理员邮箱")
		pwd      = flag.String("password", "", "管理员密码（必填）")
	)
	flag.Parse()

	if *pwd == "" {
		fmt.Println("错误：必须提供密码")
		fmt.Println("使用示例：go run scripts/create_admin.go -username=admin -email=admin@example.com -password=admin123")
		os.Exit(1)
	}

	// 初始化配置
	config.Init()

	// 初始化数据库
	database.InitDatabase()
	db := database.GetDB()

	// 创建仓库
	userRepo := repository.NewRepository(db)

	// 检查用户是否已存在
	ctx := context.Background()
	if userRepo.ExistsByUsername(ctx, *username) {
		log.Printf("用户名 '%s' 已存在，正在更新...", *username)
	}

	if userRepo.ExistsByEmail(ctx, *email) {
		log.Printf("邮箱 '%s' 已存在，正在更新...", *email)
	}

	// 加密密码
	hashedPassword, err := password.Hash(*pwd)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 创建管理员用户
	admin := &model.User{
		Username:        *username,
		Email:           *email,
		Password:        hashedPassword,
		Role:            "admin",
		IsActive:        true,
		PasswordVersion: 1,
	}

	// 检查是否已存在，如果存在则更新
	existingUser, err := userRepo.FindByUsername(ctx, *username)
	if err == nil {
		// 用户已存在，更新密码和角色
		existingUser.Password = hashedPassword
		existingUser.Role = "admin"
		existingUser.IsActive = true
		existingUser.PasswordVersion++

		if err := userRepo.Update(ctx, existingUser); err != nil {
			log.Fatalf("更新管理员账户失败: %v", err)
		}
		fmt.Printf("✅ 管理员账户已更新成功！\n")
		fmt.Printf("用户名: %s\n", existingUser.Username)
		fmt.Printf("邮箱: %s\n", existingUser.Email)
		fmt.Printf("角色: %s\n", existingUser.Role)
	} else {
		// 用户不存在，创建新用户
		if err := userRepo.Create(ctx, admin); err != nil {
			log.Fatalf("创建管理员账户失败: %v", err)
		}
		fmt.Printf("✅ 管理员账户创建成功！\n")
		fmt.Printf("用户名: %s\n", admin.Username)
		fmt.Printf("邮箱: %s\n", admin.Email)
		fmt.Printf("角色: %s\n", admin.Role)
	}

	fmt.Printf("密码: %s\n", *pwd)
	fmt.Printf("\n🎉 现在您可以使用管理员账户登录并管理其他用户了！\n")
}
