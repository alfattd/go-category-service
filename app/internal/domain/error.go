package domain

import "errors"

var (
	ErrInvalid      = errors.New("data is invalid")
	ErrNotFound     = errors.New("data not found")
	ErrDuplicate    = errors.New("data already exists")
	ErrConflict     = errors.New("conflict occurred")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrValidation   = errors.New("validation failed")
	ErrInternal     = errors.New("internal server error")
	ErrUnavailable  = errors.New("service unavailable")
)

type ErrorDetails struct {
	Code    string
	Message string
	Details map[string]interface{}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.New(message + ": " + err.Error())
}
