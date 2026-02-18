package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alfattd/category-service/internal/domain"
)

func (r *postgresCategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
	SELECT id, name, created_at, updated_at
	FROM categories
	WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var c domain.Category
	err := row.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func (r *postgresCategoryRepo) List(ctx context.Context) ([]*domain.Category, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
	SELECT id, name, created_at, updated_at
	FROM categories
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*domain.Category, 0)

	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
