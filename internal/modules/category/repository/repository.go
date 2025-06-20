package repository

import (
	"erp/internal/modules/category/model"
	"fmt"

	"gorm.io/gorm"
)

// Repository 分类仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建分类仓库
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建分类
func (r *Repository) Create(category *model.Category) error {
	// 计算层级
	if category.ParentID != nil {
		var parent model.Category
		if err := r.db.First(&parent, *category.ParentID).Error; err != nil {
			return fmt.Errorf("父分类不存在")
		}
		category.Level = parent.Level + 1
	} else {
		category.Level = 1
	}

	return r.db.Create(category).Error
}

// GetByID 根据ID获取分类
func (r *Repository) GetByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.Preload("Parent").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll 获取所有分类（分页）
func (r *Repository) GetAll(query *model.CategoryQueryRequest) ([]model.Category, int64, error) {
	var categories []model.Category
	var total int64

	db := r.db.Model(&model.Category{})

	// 添加查询条件
	if query.Name != "" {
		db = db.Where("name ILIKE ?", "%"+query.Name+"%")
	}
	if query.ParentID != nil {
		db = db.Where("parent_id = ?", *query.ParentID)
	}
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.Limit
	err := db.Preload("Parent").
		Order("level ASC, sort ASC, created_at ASC").
		Offset(offset).
		Limit(query.Limit).
		Find(&categories).Error

	return categories, total, err
}

// GetTree 获取分类树
func (r *Repository) GetTree() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("is_active = ?", true).
		Order("level ASC, sort ASC, created_at ASC").
		Find(&categories).Error
	return categories, err
}

// GetChildren 获取子分类
func (r *Repository) GetChildren(parentID uint) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("parent_id = ? AND is_active = ?", parentID, true).
		Order("sort ASC, created_at ASC").
		Find(&categories).Error
	return categories, err
}

// GetRootCategories 获取根分类
func (r *Repository) GetRootCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("parent_id IS NULL AND is_active = ?", true).
		Order("sort ASC, created_at ASC").
		Find(&categories).Error
	return categories, err
}

// Update 更新分类
func (r *Repository) Update(id uint, updates map[string]interface{}) error {
	// 如果更新了父分类，需要重新计算层级
	if parentID, ok := updates["parent_id"]; ok {
		var newLevel int
		if parentID != nil {
			// 处理 *uint 类型的 parentID
			var parentIDValue uint
			switch v := parentID.(type) {
			case *uint:
				if v != nil {
					parentIDValue = *v
				} else {
					// parentID 为 nil 指针，设置为根分类
					newLevel = 1
					updates["level"] = newLevel
					return r.db.Model(&model.Category{}).Where("id = ?", id).Updates(updates).Error
				}
			case uint:
				parentIDValue = v
			default:
				return fmt.Errorf("无效的父分类ID类型")
			}

			var parent model.Category
			if err := r.db.First(&parent, parentIDValue).Error; err != nil {
				return fmt.Errorf("父分类不存在")
			}
			newLevel = parent.Level + 1

			// 检查是否会造成循环引用
			if err := r.checkCircularReference(id, parentIDValue); err != nil {
				return err
			}
		} else {
			newLevel = 1
		}
		updates["level"] = newLevel
	}

	return r.db.Model(&model.Category{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除分类
func (r *Repository) Delete(id uint) error {
	// 检查是否有子分类
	var count int64
	if err := r.db.Model(&model.Category{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("不能删除含有子分类的分类")
	}

	// 检查是否有关联的产品
	// 这里需要根据实际的产品模型来检查
	// 假设产品表有category_id字段
	var productCount int64
	if err := r.db.Table("products").Where("category_id = ? AND deleted_at IS NULL", id).Count(&productCount).Error; err != nil {
		return err
	}
	if productCount > 0 {
		return fmt.Errorf("不能删除含有关联产品的分类")
	}

	return r.db.Delete(&model.Category{}, id).Error
}

// GetPath 获取分类路径
func (r *Repository) GetPath(id uint) ([]model.Category, error) {
	var path []model.Category
	current := id

	for current != 0 {
		var category model.Category
		if err := r.db.First(&category, current).Error; err != nil {
			return nil, err
		}
		path = append([]model.Category{category}, path...)
		if category.ParentID == nil {
			break
		}
		current = *category.ParentID
	}

	return path, nil
}

// checkCircularReference 检查循环引用
func (r *Repository) checkCircularReference(categoryID, newParentID uint) error {
	if categoryID == newParentID {
		return fmt.Errorf("分类不能设置自己为父分类")
	}

	// 检查新父分类是否是当前分类的子孙
	current := newParentID
	for current != 0 {
		var category model.Category
		if err := r.db.First(&category, current).Error; err != nil {
			return err
		}
		if category.ParentID == nil {
			break
		}
		if *category.ParentID == categoryID {
			return fmt.Errorf("不能设置子分类为父分类，这会造成循环引用")
		}
		current = *category.ParentID
	}

	return nil
}

// GetDescendants 获取所有子孙分类
func (r *Repository) GetDescendants(id uint) ([]model.Category, error) {
	var descendants []model.Category

	// 使用递归CTE查询所有子孙
	query := `
		WITH RECURSIVE category_tree AS (
			SELECT id, name, parent_id, level
			FROM categories 
			WHERE parent_id = ? AND deleted_at IS NULL
			
			UNION ALL
			
			SELECT c.id, c.name, c.parent_id, c.level
			FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
			WHERE c.deleted_at IS NULL
		)
		SELECT * FROM category_tree ORDER BY level, id
	`

	err := r.db.Raw(query, id).Scan(&descendants).Error
	return descendants, err
}

// BatchUpdateLevel 批量更新层级
func (r *Repository) BatchUpdateLevel(categoryID uint) error {
	// 获取当前分类及其所有子孙
	descendants, err := r.GetDescendants(categoryID)
	if err != nil {
		return err
	}

	// 重新计算层级
	var category model.Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		return err
	}

	// 在事务中更新层级
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, desc := range descendants {
			// 计算新层级（需要重新查询父分类的层级）
			var parent model.Category
			if err := tx.First(&parent, desc.ParentID).Error; err != nil {
				return err
			}
			newLevel := parent.Level + 1

			if err := tx.Model(&model.Category{}).
				Where("id = ?", desc.ID).
				Update("level", newLevel).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
