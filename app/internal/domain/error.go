package domain

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("data not found")
	ErrDuplicate = errors.New("data already exists")
)
