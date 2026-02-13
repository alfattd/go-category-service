package repository

import "github.com/alfattd/crud/internal/model"

type CategoryRepository interface {
	Create(*model.Category) error
	Update(*model.Category) error
	Delete(string) error
	GetByID(string) (*model.Category, error)
	List() ([]*model.Category, error)
}
