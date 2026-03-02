package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/alfattd/category-service/internal/mocks"
	"github.com/alfattd/category-service/internal/service"
	"github.com/alfattd/category-service/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestGetByID_EmptyID_ReturnsValidationError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.GetByID(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, cat)

	var valErrs *validator.ErrorsValidator
	assert.ErrorAs(t, err, &valErrs)
	assert.Contains(t, valErrs.Messages, "id is required")

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
	p := domain.PaginationParams{Page: 1, Limit: 10}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Count", mock.Anything).Return(2, nil)
	repo.On("List", mock.Anything, p).Return(categories, nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background(), p)

	assert.NoError(t, err)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 1, result.TotalPages)

	repo.AssertExpectations(t)
}

func TestList_DefaultsAppliedWhenParamsInvalid(t *testing.T) {
	p := domain.PaginationParams{Page: 0, Limit: 0}
	expectedP := domain.PaginationParams{Page: 1, Limit: 10}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Count", mock.Anything).Return(0, nil)
	repo.On("List", mock.Anything, expectedP).Return([]*domain.Category{}, nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background(), p)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Limit)

	repo.AssertExpectations(t)
}

func TestList_LimitCappedAtMax(t *testing.T) {
	p := domain.PaginationParams{Page: 1, Limit: 999}
	expectedP := domain.PaginationParams{Page: 1, Limit: 100}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Count", mock.Anything).Return(0, nil)
	repo.On("List", mock.Anything, expectedP).Return([]*domain.Category{}, nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background(), p)

	assert.NoError(t, err)
	assert.Equal(t, 100, result.Limit)

	repo.AssertExpectations(t)
}

func TestList_TotalPagesCalculation(t *testing.T) {
	p := domain.PaginationParams{Page: 1, Limit: 10}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Count", mock.Anything).Return(95, nil)
	repo.On("List", mock.Anything, p).Return([]*domain.Category{}, nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background(), p)

	assert.NoError(t, err)
	assert.Equal(t, 95, result.Total)
	assert.Equal(t, 10, result.TotalPages)

	repo.AssertExpectations(t)
}

func TestList_Empty_ReturnsEmptySliceWithMeta(t *testing.T) {
	p := domain.PaginationParams{Page: 1, Limit: 10}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Count", mock.Anything).Return(0, nil)
	repo.On("List", mock.Anything, p).Return([]*domain.Category{}, nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background(), p)

	assert.NoError(t, err)
	assert.Empty(t, result.Data)
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 1, result.TotalPages)

	repo.AssertExpectations(t)
}

func TestList_CountError_ReturnsError(t *testing.T) {
	p := domain.PaginationParams{Page: 1, Limit: 10}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Count", mock.Anything).Return(0, assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background(), p)

	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertNotCalled(t, "List")
	repo.AssertExpectations(t)
}

func TestList_RepoError_ReturnsError(t *testing.T) {
	p := domain.PaginationParams{Page: 1, Limit: 10}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Count", mock.Anything).Return(5, nil)
	repo.On("List", mock.Anything, p).Return(nil, assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	result, err := svc.List(context.Background(), p)

	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}
