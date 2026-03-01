package service

import (
	"context"
	"math"

	"github.com/alfattd/category-service/internal/domain"
)

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

func (s *CategoryService) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	if id == "" {
		return nil, domain.ErrInvalid
	}

	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context, p domain.PaginationParams) (*domain.PaginatedResult[*domain.Category], error) {
	if p.Page < 1 {
		p.Page = defaultPage
	}

	if p.Limit < 1 {
		p.Limit = defaultLimit
	} else if p.Limit > maxLimit {
		p.Limit = maxLimit
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.List(ctx, p)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &domain.PaginatedResult[*domain.Category]{
		Data:       categories,
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
