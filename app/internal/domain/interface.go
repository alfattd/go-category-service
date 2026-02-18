package domain

import "context"

type CategoryRepository interface {
	Create(ctx context.Context, c *Category) error
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
}

type CategoryService interface {
	Create(ctx context.Context, name string) (*Category, error)
	Update(ctx context.Context, id, name string) (*Category, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
}

type CategoryEventPublisher interface {
	PublishCategoryCreated(ctx context.Context, c *Category) error
	PublishCategoryUpdated(ctx context.Context, c *Category) error
	PublishCategoryDeleted(ctx context.Context, id string) error
}
