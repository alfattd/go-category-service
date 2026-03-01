package service_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/alfattd/category-service/internal/mocks"
	"github.com/alfattd/category-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// ─── Create ───────────────────────────────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)
	pub.On("PublishCategoryCreated", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), "Electronics")

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, "Electronics", cat.Name)
	assert.NotEmpty(t, cat.ID)
	assert.False(t, cat.CreatedAt.IsZero())
	assert.False(t, cat.UpdatedAt.IsZero())

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestCreate_EmptyName_ReturnsErrInvalid(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), "")

	assert.ErrorIs(t, err, domain.ErrInvalid)
	assert.Nil(t, cat)

	repo.AssertNotCalled(t, "Create")
	pub.AssertNotCalled(t, "PublishCategoryCreated")
}

func TestCreate_RepoError_ReturnsError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(domain.ErrDuplicate)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), "Electronics")

	assert.ErrorIs(t, err, domain.ErrDuplicate)
	assert.Nil(t, cat)

	pub.AssertNotCalled(t, "PublishCategoryCreated")
	repo.AssertExpectations(t)
}

func TestCreate_PublishError_StillReturnsCategory(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)
	pub.On("PublishCategoryCreated", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), "Electronics")

	assert.NoError(t, err)
	assert.NotNil(t, cat)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
	existing := &domain.Category{ID: "abc-123", Name: "Old Name"}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("GetByID", mock.Anything, "abc-123").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)
	pub.On("PublishCategoryUpdated", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "abc-123", "New Name")

	assert.NoError(t, err)
	assert.Equal(t, "New Name", cat.Name)
	assert.Equal(t, "abc-123", cat.ID)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestUpdate_EmptyID_ReturnsErrInvalid(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "", "New Name")

	assert.ErrorIs(t, err, domain.ErrInvalid)
	assert.Nil(t, cat)

	repo.AssertNotCalled(t, "GetByID")
	repo.AssertNotCalled(t, "Update")
}

func TestUpdate_EmptyName_ReturnsErrInvalid(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "abc-123", "")

	assert.ErrorIs(t, err, domain.ErrInvalid)
	assert.Nil(t, cat)

	repo.AssertNotCalled(t, "GetByID")
}

func TestUpdate_NotFound_ReturnsErrNotFound(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("GetByID", mock.Anything, "not-exist").Return(nil, domain.ErrNotFound)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "not-exist", "New Name")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Nil(t, cat)

	repo.AssertNotCalled(t, "Update")
	pub.AssertNotCalled(t, "PublishCategoryUpdated")
	repo.AssertExpectations(t)
}

func TestUpdate_RepoUpdateError_ReturnsError(t *testing.T) {
	existing := &domain.Category{ID: "abc-123", Name: "Old Name"}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("GetByID", mock.Anything, "abc-123").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "abc-123", "New Name")

	assert.Error(t, err)
	assert.Nil(t, cat)

	pub.AssertNotCalled(t, "PublishCategoryUpdated")
	repo.AssertExpectations(t)
}

func TestUpdate_PublishError_StillReturnsCategory(t *testing.T) {
	existing := &domain.Category{ID: "abc-123", Name: "Old Name"}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("GetByID", mock.Anything, "abc-123").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)
	pub.On("PublishCategoryUpdated", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "abc-123", "New Name")

	assert.NoError(t, err)
	assert.NotNil(t, cat)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Delete", mock.Anything, "abc-123").Return(nil)
	pub.On("PublishCategoryDeleted", mock.Anything, "abc-123").Return(nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	err := svc.Delete(context.Background(), "abc-123")

	assert.NoError(t, err)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestDelete_EmptyID_ReturnsErrInvalid(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	err := svc.Delete(context.Background(), "")

	assert.ErrorIs(t, err, domain.ErrInvalid)

	repo.AssertNotCalled(t, "Delete")
	pub.AssertNotCalled(t, "PublishCategoryDeleted")
}

func TestDelete_NotFound_ReturnsErrNotFound(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Delete", mock.Anything, "not-exist").Return(domain.ErrNotFound)

	svc := service.NewCategoryService(repo, pub, testLogger)
	err := svc.Delete(context.Background(), "not-exist")

	assert.ErrorIs(t, err, domain.ErrNotFound)

	pub.AssertNotCalled(t, "PublishCategoryDeleted")
	repo.AssertExpectations(t)
}

func TestDelete_PublishError_StillReturnsNil(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Delete", mock.Anything, "abc-123").Return(nil)
	pub.On("PublishCategoryDeleted", mock.Anything, "abc-123").Return(assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	err := svc.Delete(context.Background(), "abc-123")

	assert.NoError(t, err)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}
