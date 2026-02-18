package service

import (
	"context"
	"time"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/alfattd/category-service/internal/pkg/middleware"
	"github.com/google/uuid"
)

func (s *CategoryService) Create(ctx context.Context, name string) (*domain.Category, error) {
	if name == "" {
		return nil, domain.ErrInvalid
	}

	now := time.Now()

	category := &domain.Category{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	if err := s.publisher.PublishCategoryCreated(ctx, category); err != nil {
		s.log.Error("failed to publish category_created event",
			"error", err,
			"id", category.ID,
			"request_id", middleware.RequestIDFromContext(ctx),
		)
	}

	return category, nil
}

func (s *CategoryService) Update(ctx context.Context, id, name string) (*domain.Category, error) {
	if id == "" || name == "" {
		return nil, domain.ErrInvalid
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = name
	category.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	if err := s.publisher.PublishCategoryUpdated(ctx, category); err != nil {
		s.log.Error("failed to publish category_updated event",
			"error", err,
			"id", category.ID,
			"request_id", middleware.RequestIDFromContext(ctx),
		)
	}

	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalid
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.publisher.PublishCategoryDeleted(ctx, id); err != nil {
		s.log.Error("failed to publish category_deleted event",
			"error", err,
			"id", id,
			"request_id", middleware.RequestIDFromContext(ctx),
		)
	}

	return nil
}
