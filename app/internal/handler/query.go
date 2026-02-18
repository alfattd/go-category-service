package handler

import (
	"net/http"
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
	categories, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	result := make([]categoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, toCategoryResponse(c))
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Data: result,
	})
}
