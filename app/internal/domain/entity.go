package domain

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type Category struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MockCategoryRepository struct {
	mock.Mock
}
