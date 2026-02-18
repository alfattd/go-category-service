package handler

import (
	"encoding/json"
	"net/http"
)

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid request body"})
		return
	}

	category, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Message: "category created",
		Data:    toCategoryResponse(category),
	})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req updateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid request body"})
		return
	}

	category, err := h.service.Update(r.Context(), id, req.Name)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Message: "category updated",
		Data:    toCategoryResponse(category),
	})
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Message: "category deleted",
	})
}
