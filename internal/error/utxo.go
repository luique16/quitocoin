package errorpkg

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrNegativeAmount      = errors.New("amount must be positive")
	ErrUTXONotFound        = errors.New("utxo not found")
	ErrSelfTransfer        = errors.New("cannot transfer to yourself")
)
