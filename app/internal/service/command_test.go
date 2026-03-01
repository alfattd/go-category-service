package service_test

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/alfattd/category-service/internal/mocks"
	"github.com/alfattd/category-service/internal/service"
	"github.com/alfattd/category-service/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestCreate_TrimsWhitespace(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)
	pub.On("PublishCategoryCreated", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), "  Electronics  ")

	assert.NoError(t, err)
	require.NotNil(t, cat)
	assert.Equal(t, "Electronics", cat.Name, "name should be trimmed before saving")

	repo.AssertExpectations(t)
}

func TestCreate_EmptyName_ReturnsValidationError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, cat)

	var valErrs *validator.ErrorsValidator
	assert.ErrorAs(t, err, &valErrs)
	assert.Contains(t, valErrs.Messages, "name is required")

	repo.AssertNotCalled(t, "Create")
}

func TestCreate_WhitespaceOnlyName_ReturnsValidationError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), "   ")

	require.Error(t, err)
	assert.Nil(t, cat)

	var valErrs *validator.ErrorsValidator
	assert.ErrorAs(t, err, &valErrs)

	repo.AssertNotCalled(t, "Create")
}

func TestCreate_NameTooLong_ReturnsValidationError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), strings.Repeat("a", 21))

	require.Error(t, err)
	assert.Nil(t, cat)

	var valErrs *validator.ErrorsValidator
	assert.ErrorAs(t, err, &valErrs)
	assert.Contains(t, valErrs.Messages, "name must not exceed 20 characters")

	repo.AssertNotCalled(t, "Create")
}

func TestCreate_ForbiddenCharacters_ReturnsValidationError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Create(context.Background(), "<script>")

	require.Error(t, err)
	assert.Nil(t, cat)

	var valErrs *validator.ErrorsValidator
	assert.ErrorAs(t, err, &valErrs)
	assert.True(t, valErrs.HasErrors())

	repo.AssertNotCalled(t, "Create")
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
}

func TestUpdate_TrimsWhitespace(t *testing.T) {
	existing := &domain.Category{ID: "abc-123", Name: "Old Name"}

	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("GetByID", mock.Anything, "abc-123").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)
	pub.On("PublishCategoryUpdated", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "abc-123", "  New Name  ")

	assert.NoError(t, err)
	require.NotNil(t, cat)
	assert.Equal(t, "New Name", cat.Name)
}

func TestUpdate_EmptyID_ReturnsValidationError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "", "New Name")

	require.Error(t, err)
	assert.Nil(t, cat)

	var valErrs *validator.ErrorsValidator
	assert.ErrorAs(t, err, &valErrs)

	repo.AssertNotCalled(t, "GetByID")
}

func TestUpdate_InvalidName_ReturnsValidationError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	cat, err := svc.Update(context.Background(), "abc-123", "<bad>")

	require.Error(t, err)
	assert.Nil(t, cat)

	var valErrs *validator.ErrorsValidator
	assert.ErrorAs(t, err, &valErrs)

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
}

func TestDelete_EmptyID_ReturnsValidationError(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	svc := service.NewCategoryService(repo, pub, testLogger)
	err := svc.Delete(context.Background(), "")

	require.Error(t, err)

	var valErrs *validator.ErrorsValidator
	assert.ErrorAs(t, err, &valErrs)

	repo.AssertNotCalled(t, "Delete")
}

func TestDelete_NotFound_ReturnsErrNotFound(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Delete", mock.Anything, "not-exist").Return(domain.ErrNotFound)

	svc := service.NewCategoryService(repo, pub, testLogger)
	err := svc.Delete(context.Background(), "not-exist")

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestDelete_PublishError_StillReturnsNil(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	pub := new(mocks.MockCategoryEventPublisher)

	repo.On("Delete", mock.Anything, "abc-123").Return(nil)
	pub.On("PublishCategoryDeleted", mock.Anything, "abc-123").Return(assert.AnError)

	svc := service.NewCategoryService(repo, pub, testLogger)
	err := svc.Delete(context.Background(), "abc-123")

	assert.NoError(t, err)
}
