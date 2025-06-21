package service

import (
	"errors"
	"fmt"
	"reflect"

	"erp/internal/modules/attribute/model"
	"erp/internal/modules/attribute/repository"
)

// Service 属性服务
type Service struct {
	repo *repository.Repository
}

// NewService 创建属性服务
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// 属性管理相关方法

// CreateAttribute 创建属性
func (s *Service) CreateAttribute(req *model.CreateAttributeRequest) (*model.AttributeResponse, error) {
	// 检查属性名称是否已存在
	exists, err := s.repo.CheckAttributeNameExists(req.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("检查属性名称失败: %v", err)
	}
	if exists {
		return nil, errors.New("属性名称已存在")
	}

	// 创建属性对象
	attr := &model.Attribute{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		Type:         req.Type,
		Unit:         req.Unit,
		DefaultValue: req.DefaultValue,
		Sort:         req.Sort,
	}

	// 设置布尔字段的默认值
	if req.IsRequired != nil {
		attr.IsRequired = *req.IsRequired
	}
	if req.IsActive != nil {
		attr.IsActive = *req.IsActive
	} else {
		attr.IsActive = true // 默认启用
	}

	// 设置选项
	if len(req.Options) > 0 {
		if err := attr.SetOptions(req.Options); err != nil {
			return nil, fmt.Errorf("设置属性选项失败: %v", err)
		}
	}

	// 设置验证规则
	if !reflect.DeepEqual(req.Validation, model.ValidationRule{}) {
		if err := attr.SetValidation(req.Validation); err != nil {
			return nil, fmt.Errorf("设置验证规则失败: %v", err)
		}
	}

	// 保存到数据库
	if err := s.repo.CreateAttribute(attr); err != nil {
		return nil, fmt.Errorf("创建属性失败: %v", err)
	}

	// 转换为响应格式
	return s.attributeToResponse(attr)
}

// GetAttributeByID 根据ID获取属性
func (s *Service) GetAttributeByID(id uint) (*model.AttributeResponse, error) {
	attr, err := s.repo.GetAttributeByID(id)
	if err != nil {
		return nil, err
	}

	return s.attributeToResponse(attr)
}

// UpdateAttribute 更新属性
func (s *Service) UpdateAttribute(id uint, req *model.UpdateAttributeRequest) (*model.AttributeResponse, error) {
	// 获取现有属性
	attr, err := s.repo.GetAttributeByID(id)
	if err != nil {
		return nil, err
	}

	// 检查名称是否重复
	if req.Name != "" && req.Name != attr.Name {
		exists, err := s.repo.CheckAttributeNameExists(req.Name, id)
		if err != nil {
			return nil, fmt.Errorf("检查属性名称失败: %v", err)
		}
		if exists {
			return nil, errors.New("属性名称已存在")
		}
		attr.Name = req.Name
	}

	// 更新字段
	if req.DisplayName != "" {
		attr.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		attr.Description = req.Description
	}
	if req.Type != "" {
		attr.Type = req.Type
	}
	if req.Unit != "" {
		attr.Unit = req.Unit
	}
	if req.DefaultValue != "" {
		attr.DefaultValue = req.DefaultValue
	}

	if req.Sort > 0 {
		attr.Sort = req.Sort
	}

	// 更新布尔字段
	if req.IsRequired != nil {
		attr.IsRequired = *req.IsRequired
	}
	if req.IsActive != nil {
		attr.IsActive = *req.IsActive
	}

	// 更新选项
	if len(req.Options) > 0 {
		if err := attr.SetOptions(req.Options); err != nil {
			return nil, fmt.Errorf("设置属性选项失败: %v", err)
		}
	}

	// 更新验证规则
	if !reflect.DeepEqual(req.Validation, model.ValidationRule{}) {
		if err := attr.SetValidation(req.Validation); err != nil {
			return nil, fmt.Errorf("设置验证规则失败: %v", err)
		}
	}

	// 保存更新
	if err := s.repo.UpdateAttribute(attr); err != nil {
		return nil, fmt.Errorf("更新属性失败: %v", err)
	}

	return s.attributeToResponse(attr)
}

// DeleteAttribute 删除属性
func (s *Service) DeleteAttribute(id uint) error {
	// 检查是否存在属性值
	values, err := s.repo.GetAttributeValuesByAttribute(id)
	if err != nil {
		return fmt.Errorf("检查属性值失败: %v", err)
	}

	if len(values) > 0 {
		return errors.New("属性已被使用，无法删除")
	}

	return s.repo.DeleteAttribute(id)
}

// GetAttributesList 获取属性列表
func (s *Service) GetAttributesList(req *model.AttributeQueryRequest) (*model.AttributeListResponse, error) {
	attributes, total, err := s.repo.GetAttributesList(req)
	if err != nil {
		return nil, fmt.Errorf("获取属性列表失败: %v", err)
	}

	// 转换为响应格式
	var responses []model.AttributeResponse
	for _, attr := range attributes {
		resp, err := s.attributeToResponse(&attr)
		if err != nil {
			return nil, fmt.Errorf("转换属性响应失败: %v", err)
		}
		responses = append(responses, *resp)
	}

	return &model.AttributeListResponse{
		Attributes: responses,
		Pagination: model.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
		},
	}, nil
}

// GetCategoryAttributes 获取分类的属性列表
func (s *Service) GetCategoryAttributes(categoryID uint) (*model.CategoryAttributesResponse, error) {
	categoryAttributes, err := s.repo.GetAttributesByCategoryID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("获取分类属性失败: %v", err)
	}

	var responses []model.CategoryAttributeResponse
	for _, catAttr := range categoryAttributes {
		attrResp, err := s.attributeToResponse(&catAttr.Attribute)
		if err != nil {
			return nil, fmt.Errorf("转换属性响应失败: %v", err)
		}

		resp := model.CategoryAttributeResponse{
			ID:          catAttr.ID,
			CategoryID:  catAttr.CategoryID,
			AttributeID: catAttr.AttributeID,
			IsRequired:  catAttr.IsRequired,
			Sort:        catAttr.Sort,
			Attribute:   *attrResp,
			CreatedAt:   catAttr.CreatedAt,
			UpdatedAt:   catAttr.UpdatedAt,
		}
		responses = append(responses, resp)
	}

	return &model.CategoryAttributesResponse{
		CategoryID: categoryID,
		Attributes: responses,
	}, nil
}

// GetCategoryAttributesWithInheritance 获取分类的属性列表（包括继承）
func (s *Service) GetCategoryAttributesWithInheritance(categoryID uint) (*model.CategoryAttributesWithInheritanceResponse, error) {
	categoryAttributes, err := s.repo.GetCategoryAttributesWithInheritance(categoryID)
	if err != nil {
		return nil, fmt.Errorf("获取分类继承属性失败: %v", err)
	}

	var responses []model.CategoryAttributeWithInheritanceResponse
	for _, catAttr := range categoryAttributes {
		attrResp, err := s.attributeToResponse(&catAttr.Attribute)
		if err != nil {
			return nil, fmt.Errorf("转换属性响应失败: %v", err)
		}

		// 判断是否为继承属性
		isInherited := catAttr.CategoryID != categoryID

		resp := model.CategoryAttributeWithInheritanceResponse{
			ID:            catAttr.ID,
			CategoryID:    catAttr.CategoryID,
			AttributeID:   catAttr.AttributeID,
			IsRequired:    catAttr.IsRequired,
			Sort:          catAttr.Sort,
			IsInherited:   isInherited,
			InheritedFrom: nil, // 稍后填充
			Attribute:     *attrResp,
			CreatedAt:     catAttr.CreatedAt,
			UpdatedAt:     catAttr.UpdatedAt,
		}

		// 如果是继承属性，获取继承来源分类信息
		if isInherited {
			// 这里可以查询分类信息，但为了避免N+1问题，先简化处理
			inheritedFromID := catAttr.CategoryID
			resp.InheritedFrom = &inheritedFromID
		}

		responses = append(responses, resp)
	}

	return &model.CategoryAttributesWithInheritanceResponse{
		CategoryID: categoryID,
		Attributes: responses,
	}, nil
}

// GetAttributeInheritancePath 获取属性的继承路径
func (s *Service) GetAttributeInheritancePath(categoryID, attributeID uint) (*model.AttributeInheritancePathResponse, error) {
	categoryAttributes, err := s.repo.GetAttributeInheritancePath(categoryID, attributeID)
	if err != nil {
		return nil, fmt.Errorf("获取属性继承路径失败: %v", err)
	}

	var responses []model.CategoryAttributeWithInheritanceResponse
	for _, catAttr := range categoryAttributes {
		attrResp, err := s.attributeToResponse(&catAttr.Attribute)
		if err != nil {
			return nil, fmt.Errorf("转换属性响应失败: %v", err)
		}

		isInherited := catAttr.CategoryID != categoryID
		resp := model.CategoryAttributeWithInheritanceResponse{
			ID:            catAttr.ID,
			CategoryID:    catAttr.CategoryID,
			AttributeID:   catAttr.AttributeID,
			IsRequired:    catAttr.IsRequired,
			Sort:          catAttr.Sort,
			IsInherited:   isInherited,
			InheritedFrom: nil,
			Attribute:     *attrResp,
			CreatedAt:     catAttr.CreatedAt,
			UpdatedAt:     catAttr.UpdatedAt,
		}

		if isInherited {
			inheritedFromID := catAttr.CategoryID
			resp.InheritedFrom = &inheritedFromID
		}

		responses = append(responses, resp)
	}

	return &model.AttributeInheritancePathResponse{
		CategoryID:  categoryID,
		AttributeID: attributeID,
		Path:        responses,
	}, nil
}

// BindAttributeToCategory 绑定属性到分类
func (s *Service) BindAttributeToCategory(categoryID, attributeID uint, isRequired bool, sort int) (*model.CategoryAttributeResponse, error) {
	// 检查属性是否存在
	exists, err := s.repo.CheckAttributeExists(attributeID)
	if err != nil {
		return nil, fmt.Errorf("检查属性存在性失败: %v", err)
	}
	if !exists {
		return nil, errors.New("属性不存在")
	}

	// 检查是否已经绑定
	alreadyBound, err := s.repo.CheckCategoryAttributeExists(categoryID, attributeID)
	if err != nil {
		return nil, fmt.Errorf("检查分类属性关联失败: %v", err)
	}
	if alreadyBound {
		return nil, errors.New("属性已绑定到此分类")
	}

	// 绑定属性到分类
	if err := s.repo.BindAttributeToCategory(categoryID, attributeID, isRequired, sort); err != nil {
		return nil, fmt.Errorf("绑定属性到分类失败: %v", err)
	}

	// 获取绑定后的关联信息
	categoryAttr, err := s.repo.GetCategoryAttribute(categoryID, attributeID)
	if err != nil {
		return nil, fmt.Errorf("获取分类属性关联失败: %v", err)
	}

	attrResp, err := s.attributeToResponse(&categoryAttr.Attribute)
	if err != nil {
		return nil, fmt.Errorf("转换属性响应失败: %v", err)
	}

	return &model.CategoryAttributeResponse{
		ID:          categoryAttr.ID,
		CategoryID:  categoryAttr.CategoryID,
		AttributeID: categoryAttr.AttributeID,
		IsRequired:  categoryAttr.IsRequired,
		Sort:        categoryAttr.Sort,
		Attribute:   *attrResp,
		CreatedAt:   categoryAttr.CreatedAt,
		UpdatedAt:   categoryAttr.UpdatedAt,
	}, nil
}

// UnbindAttributeFromCategory 从分类解绑属性
func (s *Service) UnbindAttributeFromCategory(categoryID, attributeID uint) error {
	// 检查关联是否存在
	exists, err := s.repo.CheckCategoryAttributeExists(categoryID, attributeID)
	if err != nil {
		return fmt.Errorf("检查分类属性关联失败: %v", err)
	}
	if !exists {
		return errors.New("分类属性关联不存在")
	}

	return s.repo.UnbindAttributeFromCategory(categoryID, attributeID)
}

// UpdateCategoryAttribute 更新分类属性关联
func (s *Service) UpdateCategoryAttribute(categoryID, attributeID uint, isRequired *bool, sort *int) (*model.CategoryAttributeResponse, error) {
	// 检查关联是否存在
	exists, err := s.repo.CheckCategoryAttributeExists(categoryID, attributeID)
	if err != nil {
		return nil, fmt.Errorf("检查分类属性关联失败: %v", err)
	}
	if !exists {
		return nil, errors.New("分类属性关联不存在")
	}

	// 更新关联
	if err := s.repo.UpdateCategoryAttribute(categoryID, attributeID, isRequired, sort); err != nil {
		return nil, fmt.Errorf("更新分类属性关联失败: %v", err)
	}

	// 获取更新后的关联信息
	categoryAttr, err := s.repo.GetCategoryAttribute(categoryID, attributeID)
	if err != nil {
		return nil, fmt.Errorf("获取分类属性关联失败: %v", err)
	}

	attrResp, err := s.attributeToResponse(&categoryAttr.Attribute)
	if err != nil {
		return nil, fmt.Errorf("转换属性响应失败: %v", err)
	}

	return &model.CategoryAttributeResponse{
		ID:          categoryAttr.ID,
		CategoryID:  categoryAttr.CategoryID,
		AttributeID: categoryAttr.AttributeID,
		IsRequired:  categoryAttr.IsRequired,
		Sort:        categoryAttr.Sort,
		Attribute:   *attrResp,
		CreatedAt:   categoryAttr.CreatedAt,
		UpdatedAt:   categoryAttr.UpdatedAt,
	}, nil
}

// BatchBindAttributesToCategory 批量绑定属性到分类
func (s *Service) BatchBindAttributesToCategory(categoryID uint, attributes []model.CategoryAttributeBindRequest) error {
	// 验证所有属性是否存在
	for i, attr := range attributes {
		exists, err := s.repo.CheckAttributeExists(attr.AttributeID)
		if err != nil {
			return fmt.Errorf("检查第%d个属性存在性失败: %v", i+1, err)
		}
		if !exists {
			return fmt.Errorf("第%d个属性不存在", i+1)
		}

		// 检查是否已经绑定
		alreadyBound, err := s.repo.CheckCategoryAttributeExists(categoryID, attr.AttributeID)
		if err != nil {
			return fmt.Errorf("检查第%d个属性关联失败: %v", i+1, err)
		}
		if alreadyBound {
			return fmt.Errorf("第%d个属性已绑定到此分类", i+1)
		}
	}

	// 批量绑定
	return s.repo.BatchBindAttributesToCategory(categoryID, attributes)
}

// 属性值管理相关方法

// SetAttributeValue 设置属性值
func (s *Service) SetAttributeValue(req *model.SetAttributeValueRequest) (*model.AttributeValueResponse, error) {
	// 检查属性是否存在
	attr, err := s.repo.GetAttributeByID(req.AttributeID)
	if err != nil {
		return nil, fmt.Errorf("属性不存在: %v", err)
	}

	// 验证属性值
	if err := s.validateAttributeValue(attr, req.Value); err != nil {
		return nil, fmt.Errorf("属性值验证失败: %v", err)
	}

	// 设置属性值
	if err := s.repo.SetAttributeValue(req.AttributeID, req.EntityType, req.EntityID, req.Value); err != nil {
		return nil, fmt.Errorf("设置属性值失败: %v", err)
	}

	// 获取设置后的属性值
	value, err := s.repo.GetAttributeValue(req.AttributeID, req.EntityType, req.EntityID)
	if err != nil {
		return nil, fmt.Errorf("获取属性值失败: %v", err)
	}

	return s.attributeValueToResponse(value, true)
}

// GetAttributeValuesByEntity 获取实体的属性值
func (s *Service) GetAttributeValuesByEntity(entityType string, entityID uint) (*model.EntityAttributeValuesResponse, error) {
	values, err := s.repo.GetAttributeValuesByEntity(entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("获取实体属性值失败: %v", err)
	}

	var responses []model.AttributeValueResponse
	for _, value := range values {
		resp, err := s.attributeValueToResponse(&value, true)
		if err != nil {
			return nil, fmt.Errorf("转换属性值响应失败: %v", err)
		}
		responses = append(responses, *resp)
	}

	return &model.EntityAttributeValuesResponse{
		EntityType: entityType,
		EntityID:   entityID,
		Values:     responses,
	}, nil
}

// DeleteAttributeValue 删除属性值
func (s *Service) DeleteAttributeValue(id uint) error {
	return s.repo.DeleteAttributeValue(id)
}

// BatchSetAttributeValues 批量设置属性值
func (s *Service) BatchSetAttributeValues(values []model.SetAttributeValueRequest) error {
	// 验证所有属性值
	for i, req := range values {
		attr, err := s.repo.GetAttributeByID(req.AttributeID)
		if err != nil {
			return fmt.Errorf("第%d个属性不存在: %v", i+1, err)
		}

		if err := s.validateAttributeValue(attr, req.Value); err != nil {
			return fmt.Errorf("第%d个属性值验证失败: %v", i+1, err)
		}
	}

	// 批量设置
	return s.repo.BatchSetAttributeValues(values)
}

// 辅助方法

// attributeToResponse 将属性模型转换为响应格式
func (s *Service) attributeToResponse(attr *model.Attribute) (*model.AttributeResponse, error) {
	// 获取选项
	options, err := attr.GetOptions()
	if err != nil {
		return nil, fmt.Errorf("获取属性选项失败: %v", err)
	}

	// 获取验证规则
	validation, err := attr.GetValidation()
	if err != nil {
		return nil, fmt.Errorf("获取验证规则失败: %v", err)
	}

	return &model.AttributeResponse{
		ID:           attr.ID,
		Name:         attr.Name,
		DisplayName:  attr.DisplayName,
		Description:  attr.Description,
		Type:         attr.Type,
		Unit:         attr.Unit,
		IsRequired:   attr.IsRequired,
		IsMultiple:   attr.IsMultiple(),
		DefaultValue: attr.DefaultValue,
		Options:      options,
		Validation:   validation,
		Sort:         attr.Sort,
		IsActive:     attr.IsActive,
		CreatedAt:    attr.CreatedAt,
		UpdatedAt:    attr.UpdatedAt,
	}, nil
}

// attributeValueToResponse 将属性值模型转换为响应格式
func (s *Service) attributeValueToResponse(value *model.AttributeValue, includeAttribute bool) (*model.AttributeValueResponse, error) {
	resp := &model.AttributeValueResponse{
		ID:          value.ID,
		AttributeID: value.AttributeID,
		EntityType:  value.EntityType,
		EntityID:    value.EntityID,
		Value:       value.GetValue(),
		CreatedAt:   value.CreatedAt,
		UpdatedAt:   value.UpdatedAt,
	}

	if includeAttribute {
		attrResp, err := s.attributeToResponse(&value.Attribute)
		if err != nil {
			return nil, err
		}
		resp.Attribute = attrResp
	}

	return resp, nil
}

// validateAttributeValue 验证属性值
func (s *Service) validateAttributeValue(attr *model.Attribute, value interface{}) error {
	// 获取验证规则
	validation, err := attr.GetValidation()
	if err != nil {
		return fmt.Errorf("获取验证规则失败: %v", err)
	}

	// 检查是否必填
	if attr.IsRequired && (value == nil || value == "") {
		return errors.New("该属性为必填项")
	}

	if value == nil || value == "" {
		return nil // 非必填项可以为空
	}

	// 根据属性类型验证
	switch attr.Type {
	case model.AttributeTypeText:
		return s.validateTextValue(value, validation)
	case model.AttributeTypeNumber:
		return s.validateNumberValue(value, validation)
	case model.AttributeTypeSelect, model.AttributeTypeMultiSelect:
		return s.validateSelectValue(attr, value)
	case model.AttributeTypeBoolean:
		return s.validateBooleanValue(value)
	case model.AttributeTypeEmail:
		return s.validateEmailValue(value, validation)
	case model.AttributeTypeURL:
		return s.validateURLValue(value, validation)
	}

	return nil
}

// validateTextValue 验证文本值
func (s *Service) validateTextValue(value interface{}, validation model.ValidationRule) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("值必须为字符串类型")
	}

	if validation.MinLength != nil && len(str) < *validation.MinLength {
		return fmt.Errorf("字符串长度不能少于%d", *validation.MinLength)
	}

	if validation.MaxLength != nil && len(str) > *validation.MaxLength {
		return fmt.Errorf("字符串长度不能超过%d", *validation.MaxLength)
	}

	// TODO: 添加正则表达式验证
	if validation.Pattern != "" {
		// 这里可以添加正则表达式验证逻辑
	}

	return nil
}

// validateNumberValue 验证数字值
func (s *Service) validateNumberValue(value interface{}, validation model.ValidationRule) error {
	var num float64
	switch v := value.(type) {
	case int:
		num = float64(v)
	case int32:
		num = float64(v)
	case int64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = v
	default:
		return errors.New("值必须为数字类型")
	}

	if validation.Min != nil && num < *validation.Min {
		return fmt.Errorf("数值不能小于%v", *validation.Min)
	}

	if validation.Max != nil && num > *validation.Max {
		return fmt.Errorf("数值不能大于%v", *validation.Max)
	}

	return nil
}

// validateSelectValue 验证选择值
func (s *Service) validateSelectValue(attr *model.Attribute, value interface{}) error {
	options, err := attr.GetOptions()
	if err != nil {
		return fmt.Errorf("获取属性选项失败: %v", err)
	}

	if len(options) == 0 {
		return nil // 没有选项限制
	}

	// 创建有效值映射
	validValues := make(map[string]bool)
	for _, option := range options {
		validValues[option.Value] = true
	}

	// 验证单选
	if attr.Type == model.AttributeTypeSelect {
		str, ok := value.(string)
		if !ok {
			return errors.New("单选值必须为字符串类型")
		}
		if !validValues[str] {
			return errors.New("选择的值不在有效选项中")
		}
	}

	// 验证多选
	if attr.Type == model.AttributeTypeMultiSelect {
		switch v := value.(type) {
		case []string:
			for _, item := range v {
				if !validValues[item] {
					return fmt.Errorf("选择的值'%s'不在有效选项中", item)
				}
			}
		case []interface{}:
			for _, item := range v {
				str, ok := item.(string)
				if !ok {
					return errors.New("多选值必须为字符串数组")
				}
				if !validValues[str] {
					return fmt.Errorf("选择的值'%s'不在有效选项中", str)
				}
			}
		default:
			return errors.New("多选值必须为数组类型")
		}
	}

	return nil
}

// validateBooleanValue 验证布尔值
func (s *Service) validateBooleanValue(value interface{}) error {
	_, ok := value.(bool)
	if !ok {
		return errors.New("值必须为布尔类型")
	}
	return nil
}

// validateEmailValue 验证邮箱值
func (s *Service) validateEmailValue(value interface{}, validation model.ValidationRule) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("邮箱值必须为字符串类型")
	}

	// 简单的邮箱格式验证
	if len(str) < 5 || !contains(str, "@") || !contains(str, ".") {
		return errors.New("邮箱格式不正确")
	}

	return s.validateTextValue(value, validation)
}

// validateURLValue 验证URL值
func (s *Service) validateURLValue(value interface{}, validation model.ValidationRule) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("URL值必须为字符串类型")
	}

	// 简单的URL格式验证
	if len(str) < 7 || (!hasPrefix(str, "http://") && !hasPrefix(str, "https://")) {
		return errors.New("URL格式不正确")
	}

	return s.validateTextValue(value, validation)
}

// 辅助函数
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
