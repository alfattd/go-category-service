package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/alfattd/category-service/internal/mocks"
	"github.com/alfattd/category-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── GetByID ──────────────────────────────────────────────────────────────────

func TestGetByID_Success(t *testing.T) {
	expected := &domain.Category{
		ID:        "abc-123",
		Name:      "Electronics",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("GetByID", mock.Anything, "abc-123").Return(expected, nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.GetByID(context.Background(), "abc-123")

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, cat.ID)
	assert.Equal(t, expected.Name, cat.Name)

	repo.AssertExpectations(t)
}

func TestGetByID_EmptyID_ReturnsErrInvalid(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.GetByID(context.Background(), "")

	assert.ErrorIs(t, err, domain.ErrInvalid)
	assert.Nil(t, cat)

	repo.AssertNotCalled(t, "GetByID")
}

func TestGetByID_NotFound_ReturnsErrNotFound(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("GetByID", mock.Anything, "not-exist").Return(nil, domain.ErrNotFound)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.GetByID(context.Background(), "not-exist")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Nil(t, cat)

	repo.AssertExpectations(t)
}

func TestGetByID_RepoError_ReturnsError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("GetByID", mock.Anything, "abc-123").Return(nil, assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.GetByID(context.Background(), "abc-123")

	assert.Error(t, err)
	assert.Nil(t, cat)

	repo.AssertExpectations(t)
}

// ─── List ─────────────────────────────────────────────────────────────────────

func TestList_Success(t *testing.T) {
	categories := []*domain.Category{
		{ID: "1", Name: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", Name: "Books", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("List", mock.Anything).Return(categories, nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	repo.AssertExpectations(t)
}

func TestList_Empty_ReturnsEmptySlice(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("List", mock.Anything).Return([]*domain.Category{}, nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, result)

	repo.AssertExpectations(t)
}

func TestList_RepoError_ReturnsError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("List", mock.Anything).Return(nil, assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}
