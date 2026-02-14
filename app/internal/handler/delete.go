package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alfattd/crud/internal/dto"
	"github.com/alfattd/crud/internal/platform/rabbitmq"
	"github.com/alfattd/crud/internal/repository"
)

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID required to be filled in", http.StatusBadRequest)
		return
	}

	err := h.svc.DeleteCategory(id)
	if err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event := rabbitmq.CategoryEvent{
		ID:   id,
		Type: "deleted",
	}
	if err := h.publisher.PublishCategoryEvent(event); err != nil {
		log.Println("failed to publish category event:", err)
	}

	res := dto.DeleteCategoryResponse{
		ID:      id,
		Message: "Category successfully deleted",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
