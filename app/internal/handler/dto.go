package handler

import (
	"time"

	"github.com/alfattd/category-service/internal/domain"
)

type createCategoryRequest struct {
	Name string `json:"name"`
}

type updateCategoryRequest struct {
	Name string `json:"name"`
}

type categoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type apiResponse struct {
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func toCategoryResponse(c *domain.Category) categoryResponse {
	return categoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}
