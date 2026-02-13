package repository

import "errors"

var ErrNotFound = errors.New("data not found")
var ErrDuplicate = errors.New("data already exists")
