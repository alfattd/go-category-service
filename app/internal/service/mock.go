package service

import (
	"context"
	"testing"
	"time"

	"log/slog"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) PublishCategoryCreated(ctx context.Context, c *domain.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockPublisher) PublishCategoryUpdated(ctx context.Context, c *domain.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockPublisher) PublishCategoryDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCategoryServiceCreate(t *testing.T) {
	mockRepo := new(domain.MockCategoryRepository)
	mockPublisher := new(MockPublisher)
	logger := slog.New(slog.NewJSONHandler(nil, &slog.HandlerOptions{}))

	service := NewCategoryService(mockRepo, mockPublisher, logger)

	ctx := context.Background()

	mockRepo.On("Create", mock.MatchedBy(func(ctx context.Context) bool {
		return true
	}), mock.MatchedBy(func(c *domain.Category) bool {
		return c.Name == "Electronics"
	})).Return(nil)

	mockPublisher.On("PublishCategoryCreated", ctx, mock.MatchedBy(func(c *domain.Category) bool {
		return c.Name == "Electronics"
	})).Return(nil)

	category, err := service.Create(ctx, "Electronics")

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Electronics", category.Name)
	assert.NotEmpty(t, category.ID)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestCategoryServiceCreateInvalid(t *testing.T) {
	mockRepo := new(domain.MockCategoryRepository)
	mockPublisher := new(MockPublisher)
	logger := slog.New(slog.NewJSONHandler(nil, &slog.HandlerOptions{}))

	service := NewCategoryService(mockRepo, mockPublisher, logger)
	ctx := context.Background()

	category, err := service.Create(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Equal(t, domain.ErrInvalid, err)
}

func TestCategoryServiceUpdate(t *testing.T) {
	mockRepo := new(domain.MockCategoryRepository)
	mockPublisher := new(MockPublisher)
	logger := slog.New(slog.NewJSONHandler(nil, &slog.HandlerOptions{}))

	service := NewCategoryService(mockRepo, mockPublisher, logger)
	ctx := context.Background()

	existingCategory := &domain.Category{
		ID:        "cat-123",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", ctx, "cat-123").Return(existingCategory, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(c *domain.Category) bool {
		return c.ID == "cat-123" && c.Name == "New Name"
	})).Return(nil)
	mockPublisher.On("PublishCategoryUpdated", ctx, mock.Anything).Return(nil)

	updated, err := service.Update(ctx, "cat-123", "New Name")

	assert.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	mockRepo.AssertExpectations(t)
}
