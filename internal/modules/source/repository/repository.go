package repository

import (
	"context"
	"erp/internal/modules/source/model"

	"gorm.io/gorm"
)

type SourceRepository interface {
	Create(ctx context.Context, source *model.Source) error
	Update(ctx context.Context, source *model.Source) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.Source, error)
	List(ctx context.Context, page, pageSize int) ([]model.Source, int64, error)
	FindByCode(ctx context.Context, code string) (*model.Source, error)
	ListByStatus(ctx context.Context, status int) ([]model.Source, error)
}

type sourceRepository struct {
	db *gorm.DB
}

func NewSourceRepository(db *gorm.DB) SourceRepository {
	return &sourceRepository{db: db}
}

func (r *sourceRepository) Create(ctx context.Context, source *model.Source) error {
	return r.db.WithContext(ctx).Create(source).Error
}

func (r *sourceRepository) Update(ctx context.Context, source *model.Source) error {
	return r.db.WithContext(ctx).Save(source).Error
}

func (r *sourceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Source{}, id).Error
}

func (r *sourceRepository) FindByID(ctx context.Context, id uint) (*model.Source, error) {
	var source model.Source
	err := r.db.WithContext(ctx).First(&source, id).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}

func (r *sourceRepository) List(ctx context.Context, page, pageSize int) ([]model.Source, int64, error) {
	var sources []model.Source
	var total int64

	err := r.db.WithContext(ctx).Model(&model.Source{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&sources).Error

	return sources, total, err
}

func (r *sourceRepository) FindByCode(ctx context.Context, code string) (*model.Source, error) {
	var source model.Source
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&source).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}

func (r *sourceRepository) ListByStatus(ctx context.Context, status int) ([]model.Source, error) {
	var sources []model.Source
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&sources).Error
	return sources, err
}
