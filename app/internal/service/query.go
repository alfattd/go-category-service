package service

import (
	"context"

	"github.com/alfattd/category-service/internal/domain"
)

func (s *CategoryService) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	if id == "" {
		return nil, domain.ErrInvalid
	}

	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context) ([]*domain.Category, error) {
	return s.repo.List(ctx)
}
