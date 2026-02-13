package dto

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name,omitempty"`
}

type DeleteCategoryRequest struct {
	ID string `json:"id" binding:"required"`
}

type GetCategoryRequest struct {
	ID string `json:"id" binding:"required"`
}

type ListCategoryRequest struct {
	Name string `json:"name,omitempty"`
}
