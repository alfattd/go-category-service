package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alfattd/category-service/internal/domain"
)

type CategoryHandler struct {
	service domain.CategoryService
	log     *slog.Logger
}

func NewCategoryHandler(service domain.CategoryService, log *slog.Logger) *CategoryHandler {
	return &CategoryHandler{
		service: service,
		log:     log,
	}
}

func writeJSON(w http.ResponseWriter, status int, body apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, err error, log *slog.Logger) {
	switch {
	case errors.Is(err, domain.ErrInvalid):
		log.Warn("invalid request", "error", err)
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		log.Debug("resource not found", "error", err)
		writeJSON(w, http.StatusNotFound, apiResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrDuplicate):
		log.Warn("duplicate data", "error", err)
		writeJSON(w, http.StatusConflict, apiResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUnauthorized):
		log.Warn("unauthorized access", "error", err)
		writeJSON(w, http.StatusUnauthorized, apiResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrForbidden):
		log.Warn("forbidden access", "error", err)
		writeJSON(w, http.StatusForbidden, apiResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUnavailable):
		log.Error("service unavailable", "error", err)
		writeJSON(w, http.StatusServiceUnavailable, apiResponse{Error: "service temporarily unavailable"})
	default:
		log.Error("internal server error", "error", err)
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: "internal server error"})
	}
}
