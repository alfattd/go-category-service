package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/lib/pq"
)

type postgresCategoryRepo struct {
	db *sql.DB
}

func NewPostgresCategoryRepo(db *sql.DB) domain.CategoryRepository {
	return &postgresCategoryRepo{db: db}
}

func mapPostgresError(err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return domain.ErrDuplicate
		case "23503":
			return domain.ErrInvalid
		case "23502":
			return domain.ErrInvalid
		default:
			return fmt.Errorf("postgres error %s: %w", pgErr.Code, err)
		}
	}
	return err
}
