package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alfattd/crud/internal/dto"
)

func (h *CategoryHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.ListCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resList := make([]dto.CategoryResponse, len(categories))
	for i, k := range categories {
		resList[i] = dto.CategoryResponse{
			ID:        k.ID,
			Name:      k.Name,
			CreatedAt: k.CreatedAt,
			UpdatedAt: k.UpdatedAt,
		}
	}

	res := dto.ListCategoryResponse{
		Categories: resList,
		Total:      len(resList),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
