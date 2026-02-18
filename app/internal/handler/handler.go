package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alfattd/category-service/internal/domain"
)

type CategoryHandler struct {
	service domain.CategoryService
}

func NewCategoryHandler(service domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

type createCategoryRequest struct {
	Name string `json:"name"`
}

type updateCategoryRequest struct {
	Name string `json:"name"`
}

type categoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type apiResponse struct {
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func toCategoryResponse(c *domain.Category) categoryResponse {
	return categoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}

func writeJSON(w http.ResponseWriter, status int, body apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalid):
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, apiResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrDuplicate):
		writeJSON(w, http.StatusConflict, apiResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: "internal server error"})
	}
}
