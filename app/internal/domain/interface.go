package domain

import (
	"context"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *Category) error
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
}

type CategoryEventPublisher interface {
	PublishCategoryCreated(ctx context.Context, c *Category) error
	PublishCategoryUpdated(ctx context.Context, c *Category) error
	PublishCategoryDeleted(ctx context.Context, id string) error
}

type CategoryService interface {
	Create(ctx context.Context, name string) (*Category, error)
	Update(ctx context.Context, id, name string) (*Category, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
}

func (m *MockCategoryRepository) Create(ctx context.Context, c *Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoryRepository) Update(ctx context.Context, c *Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id string) (*Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Category), args.Error(1)
}

func (m *MockCategoryRepository) List(ctx context.Context) ([]*Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Category), args.Error(1)
}
