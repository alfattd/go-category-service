package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alfattd/crud/internal/dto"
	"github.com/alfattd/crud/internal/platform/rabbitmq"
)

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.svc.CreateCategory(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event := rabbitmq.CategoryEvent{
		ID:   category.ID,
		Name: category.Name,
		Type: "created",
	}
	if err := h.publisher.PublishCategoryEvent(event); err != nil {
		log.Println("failed to publish category event:", err)
	}

	res := dto.CreateCategoryResponse{
		Category: dto.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
