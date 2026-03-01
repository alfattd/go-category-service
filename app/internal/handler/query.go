package handler

import (
	"net/http"
	"strconv"

	"github.com/alfattd/category-service/internal/domain"
)

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	category, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Data: toCategoryResponse(category),
	})
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	p := domain.PaginationParams{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 10),
	}

	result, err := h.service.List(r.Context(), p)
	if err != nil {
		writeError(w, err)
		return
	}

	data := make([]categoryResponse, 0, len(result.Data))
	for _, c := range result.Data {
		data = append(data, toCategoryResponse(c))
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Data: paginatedResponse{
			Data: data,
			Meta: paginationMeta{
				Page:       result.Page,
				Limit:      result.Limit,
				Total:      result.Total,
				TotalPages: result.TotalPages,
				HasNext:    result.Page < result.TotalPages,
				HasPrev:    result.Page > 1,
			},
		},
	})
}

func parseIntQuery(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}

	val, err := strconv.Atoi(raw)
	if err != nil || val < 1 {
		return fallback
	}

	return val
}
