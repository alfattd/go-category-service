package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/alfattd/category-service/internal/validator"
)

type CategoryHandler struct {
	service domain.CategoryService
}

func NewCategoryHandler(service domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, err error) {
	var valErrs *validator.ErrorsValidator
	if errors.As(err, &valErrs) {
		writeJSON(w, http.StatusBadRequest, apiErrorResponse{Errors: valErrs.Messages})
		return
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, apiErrorResponse{Errors: []string{err.Error()}})
	case errors.Is(err, domain.ErrDuplicate):
		writeJSON(w, http.StatusConflict, apiErrorResponse{Errors: []string{err.Error()}})
	default:
		writeJSON(w, http.StatusInternalServerError, apiErrorResponse{Errors: []string{"internal server error"}})
	}
}
