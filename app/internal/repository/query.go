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

func (r *postgresCategoryRepo) List(ctx context.Context, p domain.PaginationParams) ([]*domain.Category, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	offset := (p.Page - 1) * p.Limit

	query := `
	SELECT id, name, created_at, updated_at
	FROM categories
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, p.Limit, offset)
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

func (r *postgresCategoryRepo) Count(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categories`).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
