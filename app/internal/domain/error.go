package domain

import "errors"

var (
	ErrInvalid   = errors.New("data is invalid")
	ErrNotFound  = errors.New("data not found")
	ErrDuplicate = errors.New("data already exists")
)
