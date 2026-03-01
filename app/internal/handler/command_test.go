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
	"github.com/alfattd/category-service/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// helper: decode response body into a map
func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	return resp
}

// helper: assert errors slice is present and non-empty
func assertErrors(t *testing.T, resp map[string]any) []any {
	t.Helper()
	raw, ok := resp["errors"]
	require.True(t, ok, "expected 'errors' key in response")
	errs, ok := raw.([]any)
	require.True(t, ok, "expected 'errors' to be an array")
	assert.NotEmpty(t, errs)
	return errs
}

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
	resp := decodeBody(t, w)
	assert.Equal(t, "category created", resp["message"])
}

func TestHandlerCreate_InvalidBody_ReturnsBadRequest(t *testing.T) {
	svc := new(mocks.MockCategoryService)
	h := handler.NewCategoryHandler(svc)

	r := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString(`invalid-json`))
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := decodeBody(t, w)
	assertErrors(t, resp)
	svc.AssertNotCalled(t, "Create")
}

func TestHandlerCreate_ValidationErrors_ReturnsBadRequestWithAllErrors(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	valErrs := &validator.ErrorsValidator{}
	valErrs.Add("name must not exceed 100 characters")
	valErrs.Add("name contains invalid characters")
	svc.On("Create", mock.Anything, mock.Anything).Return(nil, valErrs)

	h := handler.NewCategoryHandler(svc)

	body := bytes.NewBufferString(`{"name":"<toolongname>"}`)
	r := httptest.NewRequest(http.MethodPost, "/categories", body)
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := decodeBody(t, w)
	errs := assertErrors(t, resp)
	assert.Len(t, errs, 2)
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
	resp := decodeBody(t, w)
	assertErrors(t, resp)
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
	resp := decodeBody(t, w)
	assertErrors(t, resp)
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
	resp := decodeBody(t, w)
	assert.Equal(t, "category updated", resp["message"])
}

func TestHandlerUpdate_ValidationErrors_ReturnsBadRequest(t *testing.T) {
	svc := new(mocks.MockCategoryService)

	valErrs := &validator.ErrorsValidator{}
	valErrs.Add("name contains invalid characters")
	svc.On("Update", mock.Anything, "abc-123", mock.Anything).Return(nil, valErrs)

	h := handler.NewCategoryHandler(svc)

	body := bytes.NewBufferString(`{"name":"<bad>"}`)
	r := httptest.NewRequest(http.MethodPut, "/categories/abc-123", body)
	r.SetPathValue("id", "abc-123")
	w := httptest.NewRecorder()

	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := decodeBody(t, w)
	assertErrors(t, resp)
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
	resp := decodeBody(t, w)
	assertErrors(t, resp)
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
	resp := decodeBody(t, w)
	assert.Equal(t, "category deleted", resp["message"])
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
	resp := decodeBody(t, w)
	assertErrors(t, resp)
}
