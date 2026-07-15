package utxo

import (
	"context"
	"sync"

	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type Repository interface {
	GetBalance(ctx context.Context, userId string) (float32, error)
	SetBalance(ctx context.Context, userId string, amount float32) error
}

type repo struct {
	mu   sync.RWMutex
	data UTXO
}

func NewRepository() *repo {
	return &repo{
		data: make(UTXO),
	}
}

func (r *repo) GetBalance(_ context.Context, userId string) (float32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	balance, ok := r.data[userId]
	if !ok {
		return 0, errorpkg.ErrUTXONotFound
	}
	return balance, nil
}

func (r *repo) SetBalance(_ context.Context, userId string, amount float32) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[userId] = amount
	return nil
}
