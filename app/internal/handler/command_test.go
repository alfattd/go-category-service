package handler_test

import (
	"bytes"
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

// ─── Create ───────────────────────────────────────────────────────────────────

func TestHandlerCreate_Success(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	cat := &domain.Category{ID: "abc-123", Name: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	svc.On("Create", mock.Anything, "Electronics").Return(cat, nil)

	h := handler.NewCategoryHandler(svc)

	body := bytes.NewBufferString(`{"name":"Electronics"}`)
	r := httptest.NewRequest(http.MethodPost, "/categories", body)
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "category created", resp["message"])

	svc.AssertExpectations(t)
}

func TestHandlerCreate_InvalidBody_ReturnsBadRequest(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString(`invalid-json`))
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestHandlerCreate_ErrInvalid_ReturnsBadRequest(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("Create", mock.Anything, "").Return(nil, domain.ErrInvalid)

	h := handler.NewCategoryHandler(svc)

	body := bytes.NewBufferString(`{"name":""}`)
	r := httptest.NewRequest(http.MethodPost, "/categories", body)
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertExpectations(t)
}

func TestHandlerCreate_ErrDuplicate_ReturnsConflict(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("Create", mock.Anything, "Electronics").Return(nil, domain.ErrDuplicate)

	h := handler.NewCategoryHandler(svc)

	body := bytes.NewBufferString(`{"name":"Electronics"}`)
	r := httptest.NewRequest(http.MethodPost, "/categories", body)
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusConflict, w.Code)
	svc.AssertExpectations(t)
}

func TestHandlerCreate_InternalError_ReturnsInternalServerError(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("Create", mock.Anything, "Electronics").Return(nil, assert.AnError)

	h := handler.NewCategoryHandler(svc)

	body := bytes.NewBufferString(`{"name":"Electronics"}`)
	r := httptest.NewRequest(http.MethodPost, "/categories", body)
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestHandlerUpdate_Success(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	cat := &domain.Category{ID: "abc-123", Name: "New Name", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	svc.On("Update", mock.Anything, "abc-123", "New Name").Return(cat, nil)

	h := handler.NewCategoryHandler(svc)

	body := bytes.NewBufferString(`{"name":"New Name"}`)
	r := httptest.NewRequest(http.MethodPut, "/categories/abc-123", body)
	r.SetPathValue("id", "abc-123")
	w := httptest.NewRecorder()

	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "category updated", resp["message"])

	svc.AssertExpectations(t)
}

func TestHandlerUpdate_InvalidBody_ReturnsBadRequest(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodPut, "/categories/abc-123", bytes.NewBufferString(`invalid-json`))
	r.SetPathValue("id", "abc-123")
	w := httptest.NewRecorder()

	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Update")
}

func TestHandlerUpdate_NotFound_Returns404(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("Update", mock.Anything, "not-exist", "New Name").Return(nil, domain.ErrNotFound)

	h := handler.NewCategoryHandler(svc)

	body := bytes.NewBufferString(`{"name":"New Name"}`)
	r := httptest.NewRequest(http.MethodPut, "/categories/not-exist", body)
	r.SetPathValue("id", "not-exist")
	w := httptest.NewRecorder()

	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	svc.AssertExpectations(t)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestHandlerDelete_Success(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("Delete", mock.Anything, "abc-123").Return(nil)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodDelete, "/categories/abc-123", nil)
	r.SetPathValue("id", "abc-123")
	w := httptest.NewRecorder()

	h.Delete(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "category deleted", resp["message"])

	svc.AssertExpectations(t)
}

func TestHandlerDelete_NotFound_Returns404(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	svc.On("Delete", mock.Anything, "not-exist").Return(domain.ErrNotFound)

	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodDelete, "/categories/not-exist", nil)
	r.SetPathValue("id", "not-exist")
	w := httptest.NewRecorder()

	h.Delete(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	svc.AssertExpectations(t)
}
