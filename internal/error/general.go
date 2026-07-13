package errorpkg

import "errors"

var (
	ErrInternal  = errors.New("internal error")
	ErrInvalidID = errors.New("invalid id")
)
