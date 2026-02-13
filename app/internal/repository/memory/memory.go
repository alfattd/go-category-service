package memory

import (
	"sync"

	"github.com/alfattd/crud/internal/model"
	"github.com/alfattd/crud/internal/repository"
)

type InMemoryCategoryRepo struct {
	mu         sync.RWMutex
	categories map[string]*model.Category
}

func NewInMemoryCategoryRepo() *InMemoryCategoryRepo {
	return &InMemoryCategoryRepo{
		categories: make(map[string]*model.Category),
	}
}

func (r *InMemoryCategoryRepo) Create(k *model.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.categories[k.ID] = k
	return nil
}

func (r *InMemoryCategoryRepo) Update(k *model.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.categories[k.ID]
	if !ok {
		return repository.ErrNotFound
	}
	r.categories[k.ID] = k
	return nil
}

func (r *InMemoryCategoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.categories[id]
	if !ok {
		return repository.ErrNotFound
	}
	delete(r.categories, id)
	return nil
}

func (r *InMemoryCategoryRepo) GetByID(id string) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	k, ok := r.categories[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return k, nil
}

func (r *InMemoryCategoryRepo) List() ([]*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*model.Category, 0, len(r.categories))
	for _, k := range r.categories {
		list = append(list, k)
	}
	return list, nil
}
