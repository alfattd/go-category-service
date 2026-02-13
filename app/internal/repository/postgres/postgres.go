package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alfattd/crud/internal/model"
	"github.com/alfattd/crud/internal/repository"
)

type postgresCategoryRepo struct {
	db *sql.DB
}

func NewPostgresCategoryRepo(db *sql.DB) repository.CategoryRepository {
	return &postgresCategoryRepo{db: db}
}

func (r *postgresCategoryRepo) Create(k *model.Category) error {
	query := `
	INSERT INTO categories (id, name, created_at, updated_at)
	VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, k.ID, k.Name, k.CreatedAt, k.UpdatedAt)
	return err
}

func (r *postgresCategoryRepo) Update(k *model.Category) error {
	k.UpdatedAt = time.Now()
	query := `
	UPDATE categories
	SET name = $2, updated_at = $3
	WHERE id = $1
	`
	res, err := r.db.Exec(query, k.ID, k.Name, k.UpdatedAt)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *postgresCategoryRepo) Delete(id string) error {
	query := `
	DELETE FROM categories WHERE id = $1
	`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *postgresCategoryRepo) GetByID(id string) (*model.Category, error) {
	query := `
	SELECT id, name, created_at, updated_at
	FROM categories
	WHERE id = $1
	`

	row := r.db.QueryRow(query, id)
	k := &model.Category{}
	err := row.Scan(&k.ID, &k.Name, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return k, nil
}

func (r *postgresCategoryRepo) List() ([]*model.Category, error) {
	query := `
	SELECT id, name, created_at, updated_at
	FROM categories
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Category
	for rows.Next() {
		k := &model.Category{}
		if err := rows.Scan(&k.ID, &k.Name, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, k)
	}

	return result, nil
}
