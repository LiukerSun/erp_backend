package service

import (
	"context"
	"erp/internal/modules/source/model"
	"erp/internal/modules/source/repository"
	"errors"
)

type SourceService interface {
	CreateSource(ctx context.Context, source *model.Source) error
	UpdateSource(ctx context.Context, source *model.Source) error
	DeleteSource(ctx context.Context, id uint) error
	GetSource(ctx context.Context, id uint) (*model.Source, error)
	ListSources(ctx context.Context, page, pageSize int) ([]model.Source, int64, error)
	ListActiveSource(ctx context.Context) ([]model.Source, error)
}

type sourceService struct {
	repo repository.SourceRepository
}

func NewSourceService(repo repository.SourceRepository) SourceService {
	return &sourceService{repo: repo}
}

func (s *sourceService) CreateSource(ctx context.Context, source *model.Source) error {
	// 检查货源编码是否已存在
	existing, err := s.repo.FindByCode(ctx, source.Code)
	if err == nil && existing != nil {
		return errors.New("货源编码已存在")
	}

	return s.repo.Create(ctx, source)
}

func (s *sourceService) UpdateSource(ctx context.Context, source *model.Source) error {
	// 检查货源是否存在
	_, err := s.repo.FindByID(ctx, source.ID)
	if err != nil {
		return errors.New("货源不存在")
	}

	return s.repo.Update(ctx, source)
}

func (s *sourceService) DeleteSource(ctx context.Context, id uint) error {
	// 检查货源是否存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("货源不存在")
	}

	return s.repo.Delete(ctx, id)
}

func (s *sourceService) GetSource(ctx context.Context, id uint) (*model.Source, error) {
	source, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("货源不存在")
	}
	return source, nil
}

func (s *sourceService) ListSources(ctx context.Context, page, pageSize int) ([]model.Source, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}

func (s *sourceService) ListActiveSource(ctx context.Context) ([]model.Source, error) {
	return s.repo.ListByStatus(ctx, 1)
}
