package repository

import (
	productModel "erp/internal/modules/product/model"
	"erp/internal/modules/tags/model"

	"gorm.io/gorm"
)

type TagsRepository struct {
	db *gorm.DB
}

func NewTagsRepository(db *gorm.DB) *TagsRepository {
	return &TagsRepository{db: db}
}

// Create 创建标签
func (r *TagsRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// GetByID 根据ID获取标签
func (r *TagsRepository) GetByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Preload("Products").First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAll 获取所有标签
func (r *TagsRepository) GetAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Preload("Products").Find(&tags).Error
	return tags, err
}

// GetEnabled 获取所有启用的标签
func (r *TagsRepository) GetEnabled() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Where("is_enabled = ?", true).Preload("Products").Find(&tags).Error
	return tags, err
}

// Update 更新标签
func (r *TagsRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *TagsRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

// GetByName 根据名称获取标签
func (r *TagsRepository) GetByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// AddProductToTag 为标签添加产品
func (r *TagsRepository) AddProductToTag(tagID, productID uint) error {
	return r.db.Exec("INSERT INTO product_tags (product_id, tag_id, created_at) VALUES ($1, $2, NOW()) ON CONFLICT (product_id, tag_id) DO UPDATE SET created_at = NOW()", productID, tagID).Error
}

// RemoveProductFromTag 从标签移除产品
func (r *TagsRepository) RemoveProductFromTag(tagID, productID uint) error {
	return r.db.Exec("DELETE FROM product_tags WHERE product_id = $1 AND tag_id = $2", productID, tagID).Error
}

// GetProductsByTag 获取标签下的所有产品
func (r *TagsRepository) GetProductsByTag(tagID uint) ([]productModel.Product, error) {
	var products []productModel.Product
	err := r.db.Table("products").
		Joins("JOIN product_tags ON products.id = product_tags.product_id").
		Where("product_tags.tag_id = ?", tagID).
		Find(&products).Error
	return products, err
}

// GetTagsByProduct 获取产品的所有标签
func (r *TagsRepository) GetTagsByProduct(productID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Table("tags").
		Joins("JOIN product_tags ON tags.id = product_tags.tag_id").
		Where("product_tags.product_id = ?", productID).
		Find(&tags).Error
	return tags, err
}
