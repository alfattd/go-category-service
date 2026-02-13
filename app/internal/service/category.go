package service

import (
	"time"

	"github.com/alfattd/crud/internal/dto"
	"github.com/alfattd/crud/internal/model"
	"github.com/alfattd/crud/internal/repository"

	"github.com/google/uuid"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(r repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) CreateCategory(req dto.CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{
		ID:        uuid.NewString(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.Create(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) UpdateCategory(id string, req dto.UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	category.UpdatedAt = time.Now()
	err = s.repo.Update(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(id string) error {
	return s.repo.Delete(id)
}

func (s *CategoryService) GetCategoryByID(id string) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) ListCategory() ([]*model.Category, error) {
	return s.repo.List()
}
