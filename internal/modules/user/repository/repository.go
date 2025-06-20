package repository

import (
	"context"

	"erp/internal/modules/user/model"

	"gorm.io/gorm"
)

// Repository 用户仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建用户仓库
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建用户
func (r *Repository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据ID查找用户
func (r *Repository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *Repository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *Repository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *Repository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *Repository) ExistsByUsername(ctx context.Context, username string) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

// ExistsByEmail 检查邮箱是否存在
func (r *Repository) ExistsByEmail(ctx context.Context, email string) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// ExistsByEmailAndNotID 检查邮箱是否被其他用户使用
func (r *Repository) ExistsByEmailAndNotID(ctx context.Context, email string, id uint) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.User{}).Where("email = ? AND id != ?", email, id).Count(&count)
	return count > 0
}

// FindWithPagination 分页查找用户
func (r *Repository) FindWithPagination(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取用户列表
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetPasswordVersion 获取用户密码版本
func (r *Repository) GetPasswordVersion(ctx context.Context, id uint) (uint, error) {
	var user model.User
	err := r.db.WithContext(ctx).Select("password_version").First(&user, id).Error
	if err != nil {
		return 0, err
	}
	return user.PasswordVersion, nil
}

// Delete 软删除用户
func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}
