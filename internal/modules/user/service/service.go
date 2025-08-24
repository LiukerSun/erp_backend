package service

import (
	"context"
	"errors"

	"erp/config"
	"erp/internal/modules/user/model"
	"erp/internal/modules/user/repository"
	"erp/pkg/auth"
	"erp/pkg/password"

	"gorm.io/gorm"
)

// Service 用户服务
type Service struct {
	repo *repository.Repository
}

// NewService 创建用户服务
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// Register 用户注册
func (s *Service) Register(ctx context.Context, req model.RegisterRequest) (*model.Response, error) {
	// 检查用户名是否已存在
	if s.repo.ExistsByUsername(ctx, req.Username) {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if s.repo.ExistsByEmail(ctx, req.Email) {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // 默认角色
		IsActive: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 返回用户信息（不包含密码）
	return &model.Response{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	// 查找用户
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, errors.New("登录失败")
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, errors.New("账户已被禁用")
	}

	// 验证密码
	if !password.Check(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌对（包含访问token和刷新token）
	tokenPair, err := auth.GenerateTokenPair(user.ID, user.Username, user.Role, user.PasswordVersion)
	if err != nil {
		return nil, errors.New("令牌生成失败")
	}

	// 返回用户信息和令牌对
	userResponse := model.Response{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return &model.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         userResponse,
	}, nil
}

// GetProfile 获取用户资料
func (s *Service) GetProfile(ctx context.Context, userID uint) (*model.Response, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("获取用户信息失败")
	}

	return &model.Response{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *Service) UpdateProfile(ctx context.Context, userID uint, req model.UpdateProfileRequest) (*model.Response, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("获取用户信息失败")
	}

	// 更新用户信息
	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		if s.repo.ExistsByEmailAndNotID(ctx, req.Email, userID) {
			return nil, errors.New("邮箱已被其他用户使用")
		}
		user.Email = req.Email
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, errors.New("更新失败")
	}

	return &model.Response{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}, nil
}

// ChangePassword 修改密码
func (s *Service) ChangePassword(ctx context.Context, userID uint, req model.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("获取用户信息失败")
	}

	// 验证原密码
	if !password.Check(req.OldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := password.Hash(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码和密码版本（使旧token失效）
	user.Password = hashedPassword
	user.PasswordVersion++ // 增加密码版本，使所有旧token失效
	if err := s.repo.Update(ctx, user); err != nil {
		return errors.New("密码更新失败")
	}

	return nil
}

// RefreshToken 刷新访问token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenResponse, error) {
	// 解析刷新token
	claims, err := auth.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新token")
	}

	// 验证是否为刷新token
	if !auth.IsRefreshToken(claims) {
		return nil, errors.New("无效的token类型")
	}

	// 查找用户
	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("获取用户信息失败")
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, errors.New("账户已被禁用")
	}

	// 验证密码版本
	if !auth.ValidateTokenPasswordVersion(claims, user.PasswordVersion) {
		return nil, errors.New("token已失效，请重新登录")
	}

	// 生成新的访问token
	newAccessToken, err := auth.GenerateAccessToken(user.ID, user.Username, user.Role, user.PasswordVersion)
	if err != nil {
		return nil, errors.New("令牌生成失败")
	}

	return &model.RefreshTokenResponse{
		AccessToken: newAccessToken,
		ExpiresIn:   int64(config.AppConfig.JWTExpireHours * 3600), // 转换为秒
	}, nil
}

// GetUsers 获取用户列表
func (s *Service) GetUsers(ctx context.Context, page, limit int) (*model.UserListResponse, error) {
	offset := (page - 1) * limit

	users, total, err := s.repo.FindWithPagination(ctx, offset, limit)
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}

	// 转换为响应格式
	var userResponses []model.Response
	for _, user := range users {
		userResponses = append(userResponses, model.Response{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		})
	}

	return &model.UserListResponse{
		Users: userResponses,
		Pagination: model.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

// AdminCreateUser 管理员创建用户
func (s *Service) AdminCreateUser(ctx context.Context, req model.AdminCreateUserRequest) (*model.Response, error) {
	// 检查用户名是否已存在
	if s.repo.ExistsByUsername(ctx, req.Username) {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if s.repo.ExistsByEmail(ctx, req.Email) {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 返回用户信息（不包含密码）
	return &model.Response{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}, nil
}

// AdminUpdateUser 管理员更新用户
func (s *Service) AdminUpdateUser(ctx context.Context, userID uint, req model.AdminUpdateUserRequest) (*model.Response, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("获取用户信息失败")
	}

	// 更新用户信息
	updated := false

	if req.Email != "" && req.Email != user.Email {
		// 检查邮箱是否已被其他用户使用
		if s.repo.ExistsByEmailAndNotID(ctx, req.Email, userID) {
			return nil, errors.New("邮箱已被其他用户使用")
		}
		user.Email = req.Email
		updated = true
	}

	if req.Role != "" && req.Role != user.Role {
		user.Role = req.Role
		updated = true
	}

	if req.IsActive != nil && *req.IsActive != user.IsActive {
		user.IsActive = *req.IsActive
		updated = true
	}

	if updated {
		if err := s.repo.Update(ctx, user); err != nil {
			return nil, errors.New("更新失败")
		}
	}

	return &model.Response{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}, nil
}

// AdminResetUserPassword 管理员重置用户密码
func (s *Service) AdminResetUserPassword(ctx context.Context, userID uint, req model.AdminResetPasswordRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("获取用户信息失败")
	}

	// 加密新密码
	hashedPassword, err := password.Hash(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码和密码版本（使旧token失效）
	user.Password = hashedPassword
	user.PasswordVersion++ // 增加密码版本，使所有旧token失效
	if err := s.repo.Update(ctx, user); err != nil {
		return errors.New("密码重置失败")
	}

	return nil
}

// AdminDeleteUser 管理员删除用户（软删除）
func (s *Service) AdminDeleteUser(ctx context.Context, userID uint) error {
	_, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("获取用户信息失败")
	}

	if err := s.repo.Delete(ctx, userID); err != nil {
		return errors.New("删除用户失败")
	}

	return nil
}
