package service

import (
	"context"
	"erp/internal/modules/category/model"
	"erp/internal/modules/category/repository"
	"errors"

	"gorm.io/gorm"
)

// Service 分类服务
type Service struct {
	repo *repository.Repository
}

// NewService 创建分类服务
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// CreateCategory 创建分类
func (s *Service) CreateCategory(ctx context.Context, req model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
	}

	// 设置默认状态
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	} else {
		category.IsActive = true
	}

	if err := s.repo.Create(category); err != nil {
		return nil, errors.New("分类创建失败")
	}

	return s.convertToResponse(category), nil
}

// GetCategory 获取分类详情
func (s *Service) GetCategory(ctx context.Context, id uint) (*model.CategoryResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, errors.New("获取分类信息失败")
	}
	return s.convertToResponse(category), nil
}

// GetCategoryWithPath 获取带路径的分类详情
func (s *Service) GetCategoryWithPath(ctx context.Context, id uint) (*model.CategoryWithPathResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, errors.New("获取分类信息失败")
	}

	path, err := s.repo.GetPath(id)
	if err != nil {
		return nil, errors.New("获取分类路径失败")
	}

	response := &model.CategoryWithPathResponse{
		CategoryResponse: *s.convertToResponse(category),
		Path:             make([]model.CategoryPathResponse, len(path)),
	}

	for i, p := range path {
		response.Path[i] = model.CategoryPathResponse{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	return response, nil
}

// GetCategories 获取分类列表
func (s *Service) GetCategories(ctx context.Context, req model.CategoryQueryRequest) (*model.CategoryListResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	categories, total, err := s.repo.GetAll(&req)
	if err != nil {
		return nil, errors.New("获取分类列表失败")
	}

	responses := make([]model.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *s.convertToResponse(&category)
	}

	return &model.CategoryListResponse{
		Categories: responses,
		Pagination: model.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
		},
	}, nil
}

// GetCategoryTree 获取分类树
func (s *Service) GetCategoryTree(ctx context.Context) (*model.CategoryTreeListResponse, error) {
	categories, err := s.repo.GetTree()
	if err != nil {
		return nil, errors.New("获取分类树失败")
	}

	// 构建树形结构
	tree := s.buildTree(categories)

	return &model.CategoryTreeListResponse{
		Categories: tree,
	}, nil
}

// GetRootCategories 获取根分类
func (s *Service) GetRootCategories(ctx context.Context) (*model.CategoryTreeListResponse, error) {
	categories, err := s.repo.GetRootCategories()
	if err != nil {
		return nil, errors.New("获取根分类失败")
	}

	tree := make([]*model.CategoryTreeResponse, len(categories))
	for i, category := range categories {
		tree[i] = s.convertToTreeResponse(&category)
	}

	return &model.CategoryTreeListResponse{
		Categories: tree,
	}, nil
}

// GetChildrenCategories 获取子分类
func (s *Service) GetChildrenCategories(ctx context.Context, parentID uint) (*model.CategoryTreeListResponse, error) {
	categories, err := s.repo.GetChildren(parentID)
	if err != nil {
		return nil, errors.New("获取子分类失败")
	}

	tree := make([]*model.CategoryTreeResponse, len(categories))
	for i, category := range categories {
		tree[i] = s.convertToTreeResponse(&category)
	}

	return &model.CategoryTreeListResponse{
		Categories: tree,
	}, nil
}

// UpdateCategory 更新分类
func (s *Service) UpdateCategory(ctx context.Context, id uint, req model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	// 检查分类是否存在
	if _, err := s.repo.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, errors.New("获取分类信息失败")
	}

	// 禁止修改父分类，以保证属性继承的一致性
	if req.ParentID != nil {
		return nil, errors.New("不允许修改分类的层级关系，如需调整请删除后重新创建")
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, errors.New("分类更新失败")
	}

	// 返回更新后的分类
	return s.GetCategory(ctx, id)
}

// MoveCategory 移动分类（已禁用）
func (s *Service) MoveCategory(ctx context.Context, id uint, req model.MoveCategoryRequest) (*model.CategoryResponse, error) {
	return nil, errors.New("不允许移动分类层级关系，以保证属性继承的一致性。如需调整请删除后重新创建")
}

// DeleteCategory 删除分类
func (s *Service) DeleteCategory(ctx context.Context, id uint) error {
	// 检查分类是否存在
	if _, err := s.repo.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return errors.New("获取分类信息失败")
	}

	if err := s.repo.Delete(id); err != nil {
		return errors.New("分类删除失败")
	}

	return nil
}

// convertToResponse 转换为响应结构
func (s *Service) convertToResponse(category *model.Category) *model.CategoryResponse {
	return &model.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category.ParentID,
		Level:       category.Level,
		Sort:        category.Sort,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// convertToTreeResponse 转换为树形响应结构
func (s *Service) convertToTreeResponse(category *model.Category) *model.CategoryTreeResponse {
	return &model.CategoryTreeResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category.ParentID,
		Level:       category.Level,
		Sort:        category.Sort,
		IsActive:    category.IsActive,
		Children:    []*model.CategoryTreeResponse{},
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// buildTree 构建树形结构
func (s *Service) buildTree(categories []model.Category) []*model.CategoryTreeResponse {
	// 创建分类映射
	categoryMap := make(map[uint]*model.CategoryTreeResponse)
	var roots []*model.CategoryTreeResponse

	// 第一遍：创建所有节点
	for _, category := range categories {
		node := s.convertToTreeResponse(&category)
		categoryMap[category.ID] = node
	}

	// 第二遍：建立父子关系
	for _, category := range categories {
		node := categoryMap[category.ID]
		if category.ParentID == nil {
			// 根节点
			roots = append(roots, node)
		} else {
			// 子节点
			if parent, exists := categoryMap[*category.ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots
}
