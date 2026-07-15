package errorpkg

import "errors"

var (
	ErrBlockExists      = errors.New("block already exists")
	ErrBlockNotFound    = errors.New("block not found")
	ErrInvalidNonce     = errors.New("nonce does not satisfy difficulty requirement")
	ErrInvalidBlockHash = errors.New("invalid block hash")
	ErrInvalidBlockIndex = errors.New("invalid block index")
)
