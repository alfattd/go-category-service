package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alfattd/category-service/internal/domain"
)

type CategoryHandler struct {
	service domain.CategoryService
}

func NewCategoryHandler(service domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
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
