package errorpkg

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("resource not found")
	ErrInvalidID = errors.New("invalid id")
)
