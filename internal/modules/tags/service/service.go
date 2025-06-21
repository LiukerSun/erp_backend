package service

import (
	productModel "erp/internal/modules/product/model"
	"erp/internal/modules/tags/model"
	"erp/internal/modules/tags/repository"
)

type TagsService struct {
	repo *repository.TagsRepository
}

func NewTagsService(repo *repository.TagsRepository) *TagsService {
	return &TagsService{repo: repo}
}

// CreateTag 创建标签
func (s *TagsService) CreateTag(tag *model.Tag) error {
	return s.repo.Create(tag)
}

// GetTagByID 根据ID获取标签
func (s *TagsService) GetTagByID(id uint) (*model.Tag, error) {
	return s.repo.GetByID(id)
}

// GetAllTags 获取所有标签
func (s *TagsService) GetAllTags() ([]model.Tag, error) {
	return s.repo.GetAll()
}

// GetEnabledTags 获取所有启用的标签
func (s *TagsService) GetEnabledTags() ([]model.Tag, error) {
	return s.repo.GetEnabled()
}

// UpdateTag 更新标签
func (s *TagsService) UpdateTag(tag *model.Tag) error {
	return s.repo.Update(tag)
}

// DeleteTag 删除标签
func (s *TagsService) DeleteTag(id uint) error {
	return s.repo.Delete(id)
}

// GetTagByName 根据名称获取标签
func (s *TagsService) GetTagByName(name string) (*model.Tag, error) {
	return s.repo.GetByName(name)
}

// AddProductToTag 为标签添加产品
func (s *TagsService) AddProductToTag(tagID, productID uint) error {
	return s.repo.AddProductToTag(tagID, productID)
}

// RemoveProductFromTag 从标签移除产品
func (s *TagsService) RemoveProductFromTag(tagID, productID uint) error {
	return s.repo.RemoveProductFromTag(tagID, productID)
}

// GetProductsByTag 获取标签下的所有产品
func (s *TagsService) GetProductsByTag(tagID uint) ([]productModel.Product, error) {
	return s.repo.GetProductsByTag(tagID)
}

// GetTagsByProduct 获取产品的所有标签
func (s *TagsService) GetTagsByProduct(productID uint) ([]model.Tag, error) {
	return s.repo.GetTagsByProduct(productID)
}
