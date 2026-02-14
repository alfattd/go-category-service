package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alfattd/crud/internal/dto"
	"github.com/alfattd/crud/internal/platform/rabbitmq"
	"github.com/alfattd/crud/internal/repository"
)

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID required to be filled in", http.StatusBadRequest)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.svc.UpdateCategory(id, req)
	if err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event := rabbitmq.CategoryEvent{
		ID:   category.ID,
		Name: category.Name,
		Type: "updated",
	}
	if err := h.publisher.PublishCategoryEvent(event); err != nil {
		log.Println("failed to publish category event:", err)
	}

	res := dto.UpdateCategoryResponse{
		Category: dto.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
