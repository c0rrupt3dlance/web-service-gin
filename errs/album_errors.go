package errs

import "errors"

var (
	ErrNotFound            = errors.New("album not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInternalServerError = errors.New("internal server error")
)
