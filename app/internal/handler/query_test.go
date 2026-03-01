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
	json.NewDecoder(w.Body).Decode(&resp)
	assert.NotEmpty(t, resp["error"])

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
	categories := []*domain.Category{
		{ID: "1", Name: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", Name: "Books", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	svc.On("List", mock.Anything).Return(categories, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	data := resp["data"].([]any)
	assert.Len(t, data, 2)

	svc.AssertExpectations(t)
}

func TestHandlerList_Empty_ReturnsEmptyArray(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("List", mock.Anything).Return([]*domain.Category{}, nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	data := resp["data"].([]any)
	assert.Empty(t, data)

	svc.AssertExpectations(t)
}

func TestHandlerList_InternalError_Returns500(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("List", mock.Anything).Return(nil, assert.AnError)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
