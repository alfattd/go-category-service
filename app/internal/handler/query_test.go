package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/alfattd/category-service/internal/handler"
	"github.com/alfattd/category-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ─── GetByID ──────────────────────────────────────────────────────────────────

func TestHandlerGetByID_Success(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	cat := &domain.Category{ID: "abc-123", Name: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	svc.On("GetByID", mock.Anything, "abc-123").Return(cat, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories/abc-123", nil)
	r.SetPathValue("id", "abc-123")
	w := httptest.NewRecorder()

	h.GetByID(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "abc-123", data["id"])
	assert.Equal(t, "Electronics", data["name"])

	svc.AssertExpectations(t)
}

func TestHandlerGetByID_NotFound_Returns404(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("GetByID", mock.Anything, "not-exist").Return(nil, domain.ErrNotFound)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories/not-exist", nil)
	r.SetPathValue("id", "not-exist")
	w := httptest.NewRecorder()

	h.GetByID(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))

	errs, ok := resp["errors"].([]any)
	require.True(t, ok, "expected 'errors' array in response")
	assert.NotEmpty(t, errs)

	svc.AssertExpectations(t)
}

func TestHandlerGetByID_InternalError_Returns500(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("GetByID", mock.Anything, "abc-123").Return(nil, assert.AnError)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories/abc-123", nil)
	r.SetPathValue("id", "abc-123")
	w := httptest.NewRecorder()

	h.GetByID(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// ─── List ─────────────────────────────────────────────────────────────────────

func TestHandlerList_Success(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	p := domain.PaginationParams{Page: 2, Limit: 10}
	result := &domain.PaginatedResult[*domain.Category]{
		Data: []*domain.Category{
			{ID: "1", Name: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "2", Name: "Books", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		Page: 2, Limit: 10, Total: 95, TotalPages: 10,
	}

	svc.On("List", mock.Anything, p).Return(result, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories?page=2&limit=10", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	data := resp["data"].(map[string]any)
	items := data["data"].([]any)
	meta := data["meta"].(map[string]any)

	assert.Len(t, items, 2)
	assert.Equal(t, float64(2), meta["page"])
	assert.Equal(t, float64(10), meta["limit"])
	assert.Equal(t, float64(95), meta["total"])
	assert.Equal(t, float64(10), meta["total_pages"])
	assert.Equal(t, true, meta["has_next"]) // page 2 < total_pages 10
	assert.Equal(t, true, meta["has_prev"]) // page 2 > 1

	svc.AssertExpectations(t)
}

func TestHandlerList_FirstPage_HasPrevFalse(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	p := domain.PaginationParams{Page: 1, Limit: 10}
	result := &domain.PaginatedResult[*domain.Category]{
		Data: []*domain.Category{},
		Page: 1, Limit: 10, Total: 25, TotalPages: 3,
	}

	svc.On("List", mock.Anything, p).Return(result, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	meta := resp["data"].(map[string]any)["meta"].(map[string]any)

	assert.Equal(t, true, meta["has_next"])
	assert.Equal(t, false, meta["has_prev"])

	svc.AssertExpectations(t)
}

func TestHandlerList_LastPage_HasNextFalse(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	p := domain.PaginationParams{Page: 3, Limit: 10}
	result := &domain.PaginatedResult[*domain.Category]{
		Data: []*domain.Category{},
		Page: 3, Limit: 10, Total: 25, TotalPages: 3,
	}

	svc.On("List", mock.Anything, p).Return(result, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories?page=3&limit=10", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	meta := resp["data"].(map[string]any)["meta"].(map[string]any)

	assert.Equal(t, false, meta["has_next"])
	assert.Equal(t, true, meta["has_prev"])

	svc.AssertExpectations(t)
}

func TestHandlerList_SinglePage_BothFalse(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	p := domain.PaginationParams{Page: 1, Limit: 10}
	result := &domain.PaginatedResult[*domain.Category]{
		Data: []*domain.Category{},
		Page: 1, Limit: 10, Total: 3, TotalPages: 1,
	}

	svc.On("List", mock.Anything, p).Return(result, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	meta := resp["data"].(map[string]any)["meta"].(map[string]any)

	assert.Equal(t, false, meta["has_next"])
	assert.Equal(t, false, meta["has_prev"])

	svc.AssertExpectations(t)
}

func TestHandlerList_DefaultParams_WhenQueryMissing(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	p := domain.PaginationParams{Page: 1, Limit: 10}
	result := &domain.PaginatedResult[*domain.Category]{
		Data: []*domain.Category{},
		Page: 1, Limit: 10, Total: 0, TotalPages: 1,
	}

	svc.On("List", mock.Anything, p).Return(result, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandlerList_DefaultParams_WhenQueryInvalid(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	p := domain.PaginationParams{Page: 1, Limit: 10}
	result := &domain.PaginatedResult[*domain.Category]{
		Data: []*domain.Category{},
		Page: 1, Limit: 10, Total: 0, TotalPages: 1,
	}

	svc.On("List", mock.Anything, p).Return(result, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories?page=abc&limit=-5", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandlerList_InternalError_Returns500(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	p := domain.PaginationParams{Page: 1, Limit: 10}

	svc.On("List", mock.Anything, p).Return(nil, assert.AnError)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
