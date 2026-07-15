package utxo

import "context"

type Repository interface {
	GetBalance(ctx context.Context, userId string) (float32, error)
	SetBalance(ctx context.Context, userId string, amount float32) error
}
