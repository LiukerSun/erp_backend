package service

import (
	"context"
	"errors"
	"fmt"

	attributeModel "erp/internal/modules/attribute/model"
	attributeService "erp/internal/modules/attribute/service"
	categoryRepo "erp/internal/modules/category/repository"
	"erp/internal/modules/product/model"
	"erp/internal/modules/product/repository"

	"gorm.io/gorm"
)

// Service 产品服务
type Service struct {
	repo             *repository.Repository
	categoryRepo     *categoryRepo.Repository
	attributeService *attributeService.Service
}

// NewService 创建产品服务
func NewService(repo *repository.Repository, categoryRepo *categoryRepo.Repository, attributeService *attributeService.Service) *Service {
	return &Service{
		repo:             repo,
		categoryRepo:     categoryRepo,
		attributeService: attributeService,
	}
}

// CreateProduct 创建产品
func (s *Service) CreateProduct(ctx context.Context, req model.CreateProductRequest) (*model.ProductResponse, error) {
	// 创建产品
	product := &model.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, errors.New("产品创建失败")
	}

	// 返回产品信息
	return &model.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}, nil
}

// GetProduct 获取产品详情
func (s *Service) GetProduct(ctx context.Context, id uint) (*model.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, errors.New("获取产品信息失败")
	}

	return &model.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}, nil
}

// UpdateProduct 更新产品
func (s *Service) UpdateProduct(ctx context.Context, id uint, req model.UpdateProductRequest) (*model.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, errors.New("获取产品信息失败")
	}

	// 更新产品信息
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, errors.New("产品更新失败")
	}

	return &model.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}, nil
}

// DeleteProduct 删除产品
func (s *Service) DeleteProduct(ctx context.Context, id uint) error {
	// 检查产品是否存在
	if !s.repo.ExistsByID(ctx, id) {
		return errors.New("产品不存在")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.New("产品删除失败")
	}

	return nil
}

// SearchProducts 搜索产品（支持筛选和分页）
func (s *Service) SearchProducts(ctx context.Context, req model.ProductQueryRequest) (*model.ProductListResponse, error) {
	// 设置默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	// 使用查询接口，如果没有筛选条件，会返回所有产品
	products, total, err := s.repo.FindWithQuery(ctx, req, offset, req.Limit)
	if err != nil {
		return nil, errors.New("获取产品列表失败")
	}

	// 转换为响应格式
	var productResponses []model.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, model.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			CategoryID: product.CategoryID,
			CreatedAt:  product.CreatedAt,
			UpdatedAt:  product.UpdatedAt,
		})
	}

	return &model.ProductListResponse{
		Products: productResponses,
		Pagination: model.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
		},
	}, nil
}

// GetProductsByCategory 根据分类获取产品
func (s *Service) GetProductsByCategory(ctx context.Context, categoryID uint) ([]model.ProductResponse, error) {
	products, err := s.repo.FindByCategory(ctx, categoryID)
	if err != nil {
		return nil, errors.New("获取分类产品失败")
	}

	// 转换为响应格式
	var productResponses []model.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, model.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			CategoryID: product.CategoryID,
			CreatedAt:  product.CreatedAt,
			UpdatedAt:  product.UpdatedAt,
		})
	}

	return productResponses, nil
}

// GetCategoryAttributeTemplate 获取分类属性模板（用于产品创建表单）
func (s *Service) GetCategoryAttributeTemplate(ctx context.Context, categoryID uint) (*model.CategoryAttributeTemplateResponse, error) {
	// 检查分类是否存在
	if _, err := s.categoryRepo.GetByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, errors.New("获取分类信息失败")
	}

	// 获取分类的所有属性（包括继承）
	categoryAttributes, err := s.attributeService.GetCategoryAttributesWithInheritance(categoryID)
	if err != nil {
		return nil, fmt.Errorf("获取分类属性失败: %v", err)
	}

	// 转换为模板格式
	var templateItems []model.CategoryAttributeTemplateItemResponse
	for _, catAttr := range categoryAttributes.Attributes {
		templateItem := model.CategoryAttributeTemplateItemResponse{
			AttributeID:   catAttr.AttributeID,
			Name:          catAttr.Attribute.Name,
			DisplayName:   catAttr.Attribute.DisplayName,
			Type:          string(catAttr.Attribute.Type),
			Unit:          catAttr.Attribute.Unit,
			IsRequired:    catAttr.IsRequired,
			DefaultValue:  catAttr.Attribute.DefaultValue,
			Options:       catAttr.Attribute.Options,
			Validation:    catAttr.Attribute.Validation,
			Sort:          catAttr.Sort,
			IsInherited:   catAttr.IsInherited,
			InheritedFrom: catAttr.InheritedFrom,
		}
		templateItems = append(templateItems, templateItem)
	}

	return &model.CategoryAttributeTemplateResponse{
		CategoryID: categoryID,
		Attributes: templateItems,
	}, nil
}

// ValidateProductAttributes 验证产品属性
func (s *Service) ValidateProductAttributes(ctx context.Context, req model.ValidateProductAttributesRequest) (*model.ValidationResult, error) {
	// 获取分类的所有属性（包括继承）
	categoryAttributes, err := s.attributeService.GetCategoryAttributesWithInheritance(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("获取分类属性失败: %v", err)
	}

	// 创建属性映射，便于查找
	attributeMap := make(map[uint]bool)    // 记录哪些属性是必填的
	allAttributeMap := make(map[uint]bool) // 记录所有有效属性
	for _, catAttr := range categoryAttributes.Attributes {
		allAttributeMap[catAttr.AttributeID] = true
		if catAttr.IsRequired {
			attributeMap[catAttr.AttributeID] = true
		}
	}

	var validationErrors []model.AttributeValidationError

	// 检查必填属性是否都有值
	providedAttributes := make(map[uint]bool)
	for _, attr := range req.Attributes {
		providedAttributes[attr.AttributeID] = true

		// 检查属性是否属于该分类
		if !allAttributeMap[attr.AttributeID] {
			validationErrors = append(validationErrors, model.AttributeValidationError{
				AttributeID: attr.AttributeID,
				Field:       "attribute_id",
				Message:     "该属性不属于指定分类",
			})
		}
	}

	// 检查缺失的必填属性
	for attrID, isRequired := range attributeMap {
		if isRequired && !providedAttributes[attrID] {
			validationErrors = append(validationErrors, model.AttributeValidationError{
				AttributeID: attrID,
				Field:       "value",
				Message:     "必填属性值不能为空",
			})
		}
	}

	// TODO: 这里还可以添加更详细的属性值验证，比如数据类型、范围等

	return &model.ValidationResult{
		IsValid: len(validationErrors) == 0,
		Errors:  validationErrors,
	}, nil
}

// CreateProductWithAttributes 创建产品（包含属性）
func (s *Service) CreateProductWithAttributes(ctx context.Context, req model.CreateProductWithAttributesRequest) (*model.ProductWithAttributesResponse, error) {
	// 1. 验证产品属性
	validateReq := model.ValidateProductAttributesRequest{
		CategoryID: req.CategoryID,
		Attributes: req.Attributes,
	}

	validationResult, err := s.ValidateProductAttributes(ctx, validateReq)
	if err != nil {
		return nil, fmt.Errorf("属性验证失败: %v", err)
	}

	if !validationResult.IsValid {
		return nil, fmt.Errorf("属性验证不通过: %v", validationResult.Errors)
	}

	// 2. 创建产品
	product := &model.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, errors.New("产品创建失败")
	}

	// 3. 保存产品属性值
	for _, attr := range req.Attributes {
		setValueReq := &attributeModel.SetAttributeValueRequest{
			AttributeID: attr.AttributeID,
			EntityType:  "product",
			EntityID:    product.ID,
			Value:       attr.Value,
		}

		if _, err := s.attributeService.SetAttributeValue(setValueReq); err != nil {
			// 如果属性值保存失败，回滚产品创建
			s.repo.Delete(ctx, product.ID)
			return nil, fmt.Errorf("保存产品属性失败: %v", err)
		}
	}

	// 4. 返回创建的产品及其属性
	return s.GetProductWithAttributes(ctx, product.ID)
}

// GetProductWithAttributes 获取产品详情（包含属性）
func (s *Service) GetProductWithAttributes(ctx context.Context, id uint) (*model.ProductWithAttributesResponse, error) {
	// 获取产品基本信息
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, errors.New("获取产品信息失败")
	}

	// 获取产品属性值
	attributeValues, err := s.attributeService.GetAttributeValuesByEntity("product", id)
	if err != nil {
		return nil, fmt.Errorf("获取产品属性值失败: %v", err)
	}

	// 获取分类属性信息（用于判断继承）
	categoryAttributes, err := s.attributeService.GetCategoryAttributesWithInheritance(product.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("获取分类属性信息失败: %v", err)
	}

	// 创建属性继承信息映射
	inheritanceMap := make(map[uint]struct {
		isInherited   bool
		inheritedFrom *uint
		isRequired    bool
	})

	for _, catAttr := range categoryAttributes.Attributes {
		inheritanceMap[catAttr.AttributeID] = struct {
			isInherited   bool
			inheritedFrom *uint
			isRequired    bool
		}{
			isInherited:   catAttr.IsInherited,
			inheritedFrom: catAttr.InheritedFrom,
			isRequired:    catAttr.IsRequired,
		}
	}

	// 转换属性值为响应格式
	var productAttributes []model.ProductAttributeResponse
	for _, attrValue := range attributeValues.Values {
		inheritance := inheritanceMap[attrValue.AttributeID]

		productAttr := model.ProductAttributeResponse{
			AttributeID:   attrValue.AttributeID,
			AttributeName: attrValue.Attribute.Name,
			DisplayName:   attrValue.Attribute.DisplayName,
			AttributeType: string(attrValue.Attribute.Type),
			Value:         attrValue.Value,
			IsRequired:    inheritance.isRequired,
			IsInherited:   inheritance.isInherited,
			InheritedFrom: inheritance.inheritedFrom,
		}
		productAttributes = append(productAttributes, productAttr)
	}

	return &model.ProductWithAttributesResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		Attributes: productAttributes,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}, nil
}

// UpdateProductWithAttributes 更新产品（包含属性）
func (s *Service) UpdateProductWithAttributes(ctx context.Context, id uint, req model.UpdateProductWithAttributesRequest) (*model.ProductWithAttributesResponse, error) {
	// 检查产品是否存在
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, errors.New("获取产品信息失败")
	}

	// 如果要更新属性，先验证
	if len(req.Attributes) > 0 {
		categoryID := req.CategoryID
		if categoryID == 0 {
			categoryID = product.CategoryID // 使用原有分类
		}

		validateReq := model.ValidateProductAttributesRequest{
			CategoryID: categoryID,
			Attributes: req.Attributes,
		}

		validationResult, err := s.ValidateProductAttributes(ctx, validateReq)
		if err != nil {
			return nil, fmt.Errorf("属性验证失败: %v", err)
		}

		if !validationResult.IsValid {
			return nil, fmt.Errorf("属性验证不通过: %v", validationResult.Errors)
		}
	}

	// 更新产品基本信息
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, errors.New("产品更新失败")
	}

	// 更新产品属性值
	if len(req.Attributes) > 0 {
		for _, attr := range req.Attributes {
			setValueReq := &attributeModel.SetAttributeValueRequest{
				AttributeID: attr.AttributeID,
				EntityType:  "product",
				EntityID:    product.ID,
				Value:       attr.Value,
			}

			if _, err := s.attributeService.SetAttributeValue(setValueReq); err != nil {
				return nil, fmt.Errorf("更新产品属性失败: %v", err)
			}
		}
	}

	// 返回更新后的产品及其属性
	return s.GetProductWithAttributes(ctx, id)
}
