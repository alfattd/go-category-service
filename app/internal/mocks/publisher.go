package mocks

import (
	"context"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockCategoryEventPublisher struct {
	mock.Mock
}

func (m *MockCategoryEventPublisher) PublishCategoryCreated(ctx context.Context, c *domain.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoryEventPublisher) PublishCategoryUpdated(ctx context.Context, c *domain.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoryEventPublisher) PublishCategoryDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
