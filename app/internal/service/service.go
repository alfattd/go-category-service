package service

import (
	"log/slog"

	"github.com/alfattd/category-service/internal/domain"
)

type CategoryService struct {
	repo      domain.CategoryRepository
	publisher domain.CategoryEventPublisher
	log       *slog.Logger
}

var _ domain.CategoryService = (*CategoryService)(nil)

func NewCategoryService(
	repo domain.CategoryRepository,
	publisher domain.CategoryEventPublisher,
	log *slog.Logger,
) *CategoryService {
	return &CategoryService{
		repo:      repo,
		publisher: publisher,
		log:       log,
	}
}
