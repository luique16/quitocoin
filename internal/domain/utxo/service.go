package utxo

import "context"

type Service interface {
	SetBalance(ctx context.Context, userId string, amount float32) error
	GetBalance(ctx context.Context, userId string) (float32, error)
	Credit(ctx context.Context, userId string, amount float32) error
	Debit(ctx context.Context, userId string, amount float32) error
}
