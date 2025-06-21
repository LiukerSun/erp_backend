package repository

import (
	"errors"
	"fmt"

	"erp/internal/modules/attribute/model"
	categoryRepo "erp/internal/modules/category/repository"

	"gorm.io/gorm"
)

// Repository 属性仓库
type Repository struct {
	db           *gorm.DB
	categoryRepo *categoryRepo.Repository
}

// NewRepository 创建属性仓库
func NewRepository(db *gorm.DB, categoryRepo *categoryRepo.Repository) *Repository {
	return &Repository{
		db:           db,
		categoryRepo: categoryRepo,
	}
}

// 属性相关操作

// CreateAttribute 创建属性
func (r *Repository) CreateAttribute(attr *model.Attribute) error {
	return r.db.Create(attr).Error
}

// GetAttributeByID 根据ID获取属性
func (r *Repository) GetAttributeByID(id uint) (*model.Attribute, error) {
	var attr model.Attribute
	err := r.db.First(&attr, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("属性不存在")
		}
		return nil, err
	}
	return &attr, nil
}

// GetAttributeByName 根据名称获取属性
func (r *Repository) GetAttributeByName(name string) (*model.Attribute, error) {
	var attr model.Attribute
	err := r.db.Where("name = ?", name).First(&attr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("属性不存在")
		}
		return nil, err
	}
	return &attr, nil
}

// UpdateAttribute 更新属性
func (r *Repository) UpdateAttribute(attr *model.Attribute) error {
	return r.db.Save(attr).Error
}

// DeleteAttribute 删除属性（软删除）
func (r *Repository) DeleteAttribute(id uint) error {
	result := r.db.Delete(&model.Attribute{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("属性不存在")
	}
	return nil
}

// GetAttributesList 获取属性列表
func (r *Repository) GetAttributesList(req *model.AttributeQueryRequest) ([]model.Attribute, int64, error) {
	var attributes []model.Attribute
	var total int64

	query := r.db.Model(&model.Attribute{})

	// 按名称模糊搜索
	if req.Name != "" {
		query = query.Where("name LIKE ? OR display_name LIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
	}

	// 按类型筛选
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 按状态筛选
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit
	if err := query.Order("sort ASC, created_at DESC").
		Offset(offset).Limit(req.Limit).
		Find(&attributes).Error; err != nil {
		return nil, 0, err
	}

	return attributes, total, nil
}

// GetAttributesByCategoryID 根据分类ID获取属性列表
func (r *Repository) GetAttributesByCategoryID(categoryID uint) ([]model.CategoryAttribute, error) {
	var categoryAttributes []model.CategoryAttribute
	err := r.db.Preload("Attribute").
		Where("category_id = ?", categoryID).
		Order("sort ASC, created_at DESC").
		Find(&categoryAttributes).Error
	return categoryAttributes, err
}

// GetAttributesByType 根据类型获取属性列表
func (r *Repository) GetAttributesByType(attrType model.AttributeType) ([]model.Attribute, error) {
	var attributes []model.Attribute
	err := r.db.Where("type = ? AND is_active = ?", attrType, true).
		Order("sort ASC, created_at DESC").
		Find(&attributes).Error
	return attributes, err
}

// 分类属性关联相关操作

// BindAttributeToCategory 绑定属性到分类
func (r *Repository) BindAttributeToCategory(categoryID, attributeID uint, isRequired bool, sort int) error {
	categoryAttr := &model.CategoryAttribute{
		CategoryID:  categoryID,
		AttributeID: attributeID,
		IsRequired:  isRequired,
		Sort:        sort,
	}
	return r.db.Create(categoryAttr).Error
}

// UnbindAttributeFromCategory 从分类解绑属性
func (r *Repository) UnbindAttributeFromCategory(categoryID, attributeID uint) error {
	result := r.db.Where("category_id = ? AND attribute_id = ?", categoryID, attributeID).
		Delete(&model.CategoryAttribute{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("分类属性关联不存在")
	}
	return nil
}

// UpdateCategoryAttribute 更新分类属性关联
func (r *Repository) UpdateCategoryAttribute(categoryID, attributeID uint, isRequired *bool, sort *int) error {
	updates := make(map[string]interface{})
	if isRequired != nil {
		updates["is_required"] = *isRequired
	}
	if sort != nil {
		updates["sort"] = *sort
	}

	if len(updates) == 0 {
		return nil // 没有需要更新的字段
	}

	result := r.db.Model(&model.CategoryAttribute{}).
		Where("category_id = ? AND attribute_id = ?", categoryID, attributeID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("分类属性关联不存在")
	}
	return nil
}

// GetCategoryAttribute 获取分类属性关联
func (r *Repository) GetCategoryAttribute(categoryID, attributeID uint) (*model.CategoryAttribute, error) {
	var categoryAttr model.CategoryAttribute
	err := r.db.Preload("Attribute").
		Where("category_id = ? AND attribute_id = ?", categoryID, attributeID).
		First(&categoryAttr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类属性关联不存在")
		}
		return nil, err
	}
	return &categoryAttr, nil
}

// BatchBindAttributesToCategory 批量绑定属性到分类
func (r *Repository) BatchBindAttributesToCategory(categoryID uint, attributes []model.CategoryAttributeBindRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, attr := range attributes {
			categoryAttr := &model.CategoryAttribute{
				CategoryID:  categoryID,
				AttributeID: attr.AttributeID,
				IsRequired:  attr.IsRequired,
				Sort:        attr.Sort,
			}
			if err := tx.Create(categoryAttr).Error; err != nil {
				return fmt.Errorf("绑定属性 %d 失败: %v", attr.AttributeID, err)
			}
		}
		return nil
	})
}

// UnbindAllAttributesFromCategory 从分类解绑所有属性
func (r *Repository) UnbindAllAttributesFromCategory(categoryID uint) error {
	return r.db.Where("category_id = ?", categoryID).Delete(&model.CategoryAttribute{}).Error
}

// CheckCategoryAttributeExists 检查分类属性关联是否存在（不包括已删除的关联）
func (r *Repository) CheckCategoryAttributeExists(categoryID, attributeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.CategoryAttribute{}).
		Where("category_id = ? AND attribute_id = ? AND deleted_at IS NULL", categoryID, attributeID).
		Count(&count).Error
	return count > 0, err
}

// GetCategoryAttributesWithInheritance 获取分类的属性（包括继承）
func (r *Repository) GetCategoryAttributesWithInheritance(categoryID uint) ([]model.CategoryAttribute, error) {
	// 使用递归CTE查询获取分类路径上的所有属性绑定
	query := `
		WITH RECURSIVE category_path AS (
			-- 起始分类
			SELECT id, parent_id, name, level
			FROM categories 
			WHERE id = ? AND deleted_at IS NULL
			
			UNION ALL
			
			-- 递归查找父分类
			SELECT c.id, c.parent_id, c.name, c.level
			FROM categories c
			INNER JOIN category_path cp ON c.id = cp.parent_id
			WHERE c.deleted_at IS NULL
		)
		SELECT ca.*, a.* 
		FROM category_attributes ca
		INNER JOIN category_path cp ON ca.category_id = cp.id
		INNER JOIN attributes a ON ca.attribute_id = a.id
		WHERE ca.deleted_at IS NULL AND a.deleted_at IS NULL
		ORDER BY cp.level DESC, ca.sort ASC, ca.created_at ASC
	`

	type CategoryAttributeWithLevel struct {
		model.CategoryAttribute
		Level int `gorm:"column:level"`
	}

	var results []CategoryAttributeWithLevel
	err := r.db.Raw(query, categoryID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 处理属性优先级：子分类的设置优先于父分类
	attributeMap := make(map[uint]*model.CategoryAttribute)
	var categoryAttributes []model.CategoryAttribute

	for _, result := range results {
		// 如果属性还没有被添加，或者当前分类层级更深（子分类优先）
		if _, exists := attributeMap[result.AttributeID]; !exists || result.CategoryID == categoryID {
			// 预加载属性信息
			if err := r.db.Preload("Attribute").First(&result.CategoryAttribute, result.ID).Error; err != nil {
				continue // 跳过出错的记录
			}

			if !exists {
				categoryAttributes = append(categoryAttributes, result.CategoryAttribute)
			} else {
				// 替换为子分类的设置
				for i, ca := range categoryAttributes {
					if ca.AttributeID == result.AttributeID {
						categoryAttributes[i] = result.CategoryAttribute
						break
					}
				}
			}
			attributeMap[result.AttributeID] = &result.CategoryAttribute
		}
	}

	return categoryAttributes, nil
}

// GetAttributeInheritancePath 获取属性在分类路径中的继承信息
func (r *Repository) GetAttributeInheritancePath(categoryID, attributeID uint) ([]model.CategoryAttribute, error) {
	query := `
		WITH RECURSIVE category_path AS (
			-- 起始分类
			SELECT id, parent_id, name, level
			FROM categories 
			WHERE id = ? AND deleted_at IS NULL
			
			UNION ALL
			
			-- 递归查找父分类
			SELECT c.id, c.parent_id, c.name, c.level
			FROM categories c
			INNER JOIN category_path cp ON c.id = cp.parent_id
			WHERE c.deleted_at IS NULL
		)
		SELECT ca.* 
		FROM category_attributes ca
		INNER JOIN category_path cp ON ca.category_id = cp.id
		WHERE ca.attribute_id = ? AND ca.deleted_at IS NULL
		ORDER BY cp.level DESC
	`

	var categoryAttributes []model.CategoryAttribute
	err := r.db.Raw(query, categoryID, attributeID).Preload("Attribute").Find(&categoryAttributes).Error
	return categoryAttributes, err
}

// 属性值相关操作

// CreateAttributeValue 创建属性值
func (r *Repository) CreateAttributeValue(value *model.AttributeValue) error {
	return r.db.Create(value).Error
}

// GetAttributeValueByID 根据ID获取属性值
func (r *Repository) GetAttributeValueByID(id uint) (*model.AttributeValue, error) {
	var value model.AttributeValue
	err := r.db.Preload("Attribute").First(&value, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("属性值不存在")
		}
		return nil, err
	}
	return &value, nil
}

// UpdateAttributeValue 更新属性值
func (r *Repository) UpdateAttributeValue(value *model.AttributeValue) error {
	return r.db.Save(value).Error
}

// DeleteAttributeValue 删除属性值（软删除）
func (r *Repository) DeleteAttributeValue(id uint) error {
	result := r.db.Delete(&model.AttributeValue{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("属性值不存在")
	}
	return nil
}

// GetAttributeValuesByEntity 根据实体获取属性值列表
func (r *Repository) GetAttributeValuesByEntity(entityType string, entityID uint) ([]model.AttributeValue, error) {
	var values []model.AttributeValue
	err := r.db.Preload("Attribute").
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Find(&values).Error
	return values, err
}

// GetAttributeValuesByAttribute 根据属性ID获取属性值列表
func (r *Repository) GetAttributeValuesByAttribute(attributeID uint) ([]model.AttributeValue, error) {
	var values []model.AttributeValue
	err := r.db.Preload("Attribute").
		Where("attribute_id = ?", attributeID).
		Find(&values).Error
	return values, err
}

// GetAttributeValue 获取特定的属性值
func (r *Repository) GetAttributeValue(attributeID uint, entityType string, entityID uint) (*model.AttributeValue, error) {
	var value model.AttributeValue
	err := r.db.Preload("Attribute").
		Where("attribute_id = ? AND entity_type = ? AND entity_id = ?", attributeID, entityType, entityID).
		First(&value).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回nil表示未找到，但不是错误
		}
		return nil, err
	}
	return &value, nil
}

// SetAttributeValue 设置属性值（如果存在则更新，不存在则创建）
func (r *Repository) SetAttributeValue(attributeID uint, entityType string, entityID uint, value interface{}) error {
	// 检查是否已存在
	existingValue, err := r.GetAttributeValue(attributeID, entityType, entityID)
	if err != nil {
		return err
	}

	if existingValue != nil {
		// 更新现有值
		if err := existingValue.SetValue(value); err != nil {
			return err
		}
		return r.UpdateAttributeValue(existingValue)
	} else {
		// 创建新值
		newValue := &model.AttributeValue{
			AttributeID: attributeID,
			EntityType:  entityType,
			EntityID:    entityID,
		}
		if err := newValue.SetValue(value); err != nil {
			return err
		}
		return r.CreateAttributeValue(newValue)
	}
}

// DeleteAttributeValuesByEntity 删除实体的所有属性值
func (r *Repository) DeleteAttributeValuesByEntity(entityType string, entityID uint) error {
	return r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Delete(&model.AttributeValue{}).Error
}

// DeleteAttributeValuesByAttribute 删除属性的所有值
func (r *Repository) DeleteAttributeValuesByAttribute(attributeID uint) error {
	return r.db.Where("attribute_id = ?", attributeID).
		Delete(&model.AttributeValue{}).Error
}

// BatchSetAttributeValues 批量设置属性值
func (r *Repository) BatchSetAttributeValues(values []model.SetAttributeValueRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, req := range values {
			if err := r.SetAttributeValue(req.AttributeID, req.EntityType, req.EntityID, req.Value); err != nil {
				return fmt.Errorf("设置属性值失败 (AttributeID: %d, EntityType: %s, EntityID: %d): %v",
					req.AttributeID, req.EntityType, req.EntityID, err)
			}
		}
		return nil
	})
}

// GetAttributeValuesWithPagination 分页获取属性值
func (r *Repository) GetAttributeValuesWithPagination(entityType string, page, limit int) ([]model.AttributeValue, int64, error) {
	var values []model.AttributeValue
	var total int64

	query := r.db.Model(&model.AttributeValue{}).Preload("Attribute")

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&values).Error; err != nil {
		return nil, 0, err
	}

	return values, total, nil
}

// CheckAttributeExists 检查属性是否存在（不包括已删除的属性）
func (r *Repository) CheckAttributeExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Attribute{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	return count > 0, err
}

// CheckAttributeNameExists 检查属性名称是否存在（用于创建时检查重复，不包括已删除的属性）
func (r *Repository) CheckAttributeNameExists(name string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Attribute{}).Where("name = ? AND deleted_at IS NULL", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// 级联更新相关方法

// CascadeBindAttributeToDescendants 级联绑定属性到所有子分类
func (r *Repository) CascadeBindAttributeToDescendants(parentCategoryID, attributeID uint, isRequired bool, sort int) error {
	// 获取所有子孙分类
	descendantIDs, err := r.categoryRepo.GetAllDescendants(parentCategoryID)
	if err != nil {
		return fmt.Errorf("获取子分类失败: %v", err)
	}

	if len(descendantIDs) == 0 {
		return nil // 没有子分类，无需级联
	}

	// 在事务中执行级联绑定
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, categoryID := range descendantIDs {
			// 检查子分类是否已经绑定了此属性
			var count int64
			err := tx.Model(&model.CategoryAttribute{}).
				Where("category_id = ? AND attribute_id = ?", categoryID, attributeID).
				Count(&count).Error
			if err != nil {
				return fmt.Errorf("检查分类%d的属性绑定失败: %v", categoryID, err)
			}

			// 如果子分类还没有绑定此属性，则继承父分类的绑定
			if count == 0 {
				categoryAttr := &model.CategoryAttribute{
					CategoryID:  categoryID,
					AttributeID: attributeID,
					IsRequired:  isRequired,
					Sort:        sort,
				}
				if err := tx.Create(categoryAttr).Error; err != nil {
					return fmt.Errorf("为分类%d绑定继承属性失败: %v", categoryID, err)
				}
			}
		}
		return nil
	})
}

// CascadeUnbindAttributeFromDescendants 级联解绑属性从所有子分类
func (r *Repository) CascadeUnbindAttributeFromDescendants(parentCategoryID, attributeID uint) error {
	// 获取所有子孙分类
	descendantIDs, err := r.categoryRepo.GetAllDescendants(parentCategoryID)
	if err != nil {
		return fmt.Errorf("获取子分类失败: %v", err)
	}

	if len(descendantIDs) == 0 {
		return nil // 没有子分类，无需级联
	}

	// 在事务中执行级联解绑
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, categoryID := range descendantIDs {
			// 检查此属性绑定是否为继承而来（即：在分类继承路径中的父级有此属性绑定）
			isInherited, err := r.isAttributeInheritedFromParent(categoryID, attributeID, parentCategoryID)
			if err != nil {
				return fmt.Errorf("检查属性继承关系失败: %v", err)
			}

			// 只有当此属性确实是继承自要解绑的父分类时，才进行解绑
			if isInherited {
				result := tx.Where("category_id = ? AND attribute_id = ?", categoryID, attributeID).
					Delete(&model.CategoryAttribute{})
				if result.Error != nil {
					return fmt.Errorf("从分类%d解绑继承属性失败: %v", categoryID, result.Error)
				}
			}
		}
		return nil
	})
}

// CascadeUpdateAttributeInDescendants 级联更新属性设置到所有子分类
func (r *Repository) CascadeUpdateAttributeInDescendants(parentCategoryID, attributeID uint, isRequired *bool, sort *int) error {
	// 获取所有子孙分类
	descendantIDs, err := r.categoryRepo.GetAllDescendants(parentCategoryID)
	if err != nil {
		return fmt.Errorf("获取子分类失败: %v", err)
	}

	if len(descendantIDs) == 0 {
		return nil // 没有子分类，无需级联
	}

	// 准备更新字段
	updates := make(map[string]interface{})
	if isRequired != nil {
		updates["is_required"] = *isRequired
	}
	if sort != nil {
		updates["sort"] = *sort
	}

	if len(updates) == 0 {
		return nil // 没有需要更新的字段
	}

	// 在事务中执行级联更新
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, categoryID := range descendantIDs {
			// 检查此属性绑定是否为继承而来
			isInherited, err := r.isAttributeInheritedFromParent(categoryID, attributeID, parentCategoryID)
			if err != nil {
				return fmt.Errorf("检查属性继承关系失败: %v", err)
			}

			// 只有当此属性确实是继承自要更新的父分类时，才进行更新
			if isInherited {
				result := tx.Model(&model.CategoryAttribute{}).
					Where("category_id = ? AND attribute_id = ?", categoryID, attributeID).
					Updates(updates)
				if result.Error != nil {
					return fmt.Errorf("更新分类%d的继承属性设置失败: %v", categoryID, result.Error)
				}
			}
		}
		return nil
	})
}

// isAttributeInheritedFromParent 检查分类的某个属性是否继承自指定的父分类
func (r *Repository) isAttributeInheritedFromParent(categoryID, attributeID, parentCategoryID uint) (bool, error) {
	// 检查当前分类是否有自己绑定的此属性
	var ownBinding model.CategoryAttribute
	err := r.db.Where("category_id = ? AND attribute_id = ?", categoryID, attributeID).
		First(&ownBinding).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // 分类没有此属性绑定
	}
	if err != nil {
		return false, err
	}

	// 如果分类有此属性绑定，需要判断是否为继承
	// 通过查询分类的继承路径来判断
	query := `
		WITH RECURSIVE category_path AS (
			-- 从当前分类开始
			SELECT id, parent_id, 1 as level
			FROM categories 
			WHERE id = ? AND deleted_at IS NULL
			
			UNION ALL
			
			-- 递归查找父分类
			SELECT c.id, c.parent_id, cp.level + 1
			FROM categories c
			INNER JOIN category_path cp ON c.id = cp.parent_id
			WHERE c.deleted_at IS NULL AND cp.level < 10 -- 防止无限递归
		)
		SELECT COUNT(*) as count
		FROM category_path cp
		INNER JOIN category_attributes ca ON ca.category_id = cp.id
		WHERE ca.attribute_id = ? 
		  AND cp.id = ?
		  AND cp.level > 1 -- 排除自己
	`

	var count int64
	err = r.db.Raw(query, categoryID, attributeID, parentCategoryID).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// RebuildInheritanceForCategory 重建分类的属性继承关系（用于修复不一致的情况）
func (r *Repository) RebuildInheritanceForCategory(categoryID uint) error {
	// 获取分类的继承路径上的所有属性绑定
	inheritedAttributes, err := r.GetCategoryAttributesWithInheritance(categoryID)
	if err != nil {
		return fmt.Errorf("获取继承属性失败: %v", err)
	}

	// 获取分类自己的属性绑定
	ownAttributes, err := r.GetAttributesByCategoryID(categoryID)
	if err != nil {
		return fmt.Errorf("获取自有属性失败: %v", err)
	}

	// 创建自有属性映射
	ownAttributeMap := make(map[uint]bool)
	for _, attr := range ownAttributes {
		ownAttributeMap[attr.AttributeID] = true
	}

	// 在事务中重建继承关系
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除所有继承的属性绑定（保留自有绑定）
		// 这个操作比较复杂，需要精确识别哪些是继承的

		// 简化处理：重新添加缺失的继承属性
		for _, inheritedAttr := range inheritedAttributes {
			// 如果是继承属性且分类自己没有绑定
			if inheritedAttr.CategoryID != categoryID && !ownAttributeMap[inheritedAttr.AttributeID] {
				// 检查是否已存在绑定
				var count int64
				err := tx.Model(&model.CategoryAttribute{}).
					Where("category_id = ? AND attribute_id = ?", categoryID, inheritedAttr.AttributeID).
					Count(&count).Error
				if err != nil {
					return err
				}

				// 如果不存在，则创建继承绑定
				if count == 0 {
					categoryAttr := &model.CategoryAttribute{
						CategoryID:  categoryID,
						AttributeID: inheritedAttr.AttributeID,
						IsRequired:  inheritedAttr.IsRequired,
						Sort:        inheritedAttr.Sort,
					}
					if err := tx.Create(categoryAttr).Error; err != nil {
						return fmt.Errorf("创建继承属性绑定失败: %v", err)
					}
				}
			}
		}
		return nil
	})
}
