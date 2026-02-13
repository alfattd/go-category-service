package dto

import "time"

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteCategoryResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type CreateCategoryResponse struct {
	Category CategoryResponse `json:"category"`
}

type UpdateCategoryResponse struct {
	Category CategoryResponse `json:"category"`
}

type GetCategoryResponse struct {
	Category CategoryResponse `json:"category"`
}

type ListCategoryResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int                `json:"total"`
}
