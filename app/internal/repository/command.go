package repository

import (
	"context"

	"github.com/alfattd/category-service/internal/domain"
)

func (r *postgresCategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
	INSERT INTO categories (id, name, created_at, updated_at)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, c.ID, c.Name, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return mapPostgresError(err)
	}

	return nil
}

func (r *postgresCategoryRepo) Update(ctx context.Context, c *domain.Category) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
	UPDATE categories
	SET name = $1,
		updated_at = $2
	WHERE id = $3
	`

	res, err := r.db.ExecContext(ctx, query, c.Name, c.UpdatedAt, c.ID)
	if err != nil {
		return mapPostgresError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *postgresCategoryRepo) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `DELETE FROM categories WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return mapPostgresError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
