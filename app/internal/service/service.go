package service

import (
	"github.com/alfattd/category-service/internal/domain"
)

type CategoryService struct {
	repo      domain.CategoryRepository
	publisher domain.CategoryEventPublisher
}

var _ domain.CategoryService = (*CategoryService)(nil)

func NewCategoryService(
	repo domain.CategoryRepository,
	publisher domain.CategoryEventPublisher,
) *CategoryService {
	return &CategoryService{
		repo:      repo,
		publisher: publisher,
	}
}
