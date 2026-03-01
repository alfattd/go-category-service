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

type paginationMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type paginatedResponse struct {
	Data []categoryResponse `json:"data"`
	Meta paginationMeta     `json:"meta"`
}

type apiResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type apiErrorResponse struct {
	Errors []string `json:"errors"`
}

func toCategoryResponse(c *domain.Category) categoryResponse {
	return categoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}
