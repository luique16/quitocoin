package utxo

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

const utxoKeyPrefix = "utxo:"

type Repository interface {
	GetBalance(ctx context.Context, userId string) (float32, error)
	SetBalance(ctx context.Context, userId string, amount float32) error
	Clear(ctx context.Context) error
}

type repo struct {
	rdb redis.Cmdable
}

func NewRepository(rdb redis.Cmdable) *repo {
	return &repo{rdb: rdb}
}

func (r *repo) GetBalance(ctx context.Context, userId string) (float32, error) {
	val, err := r.rdb.Get(ctx, utxoKeyPrefix+userId).Float32()
	if errors.Is(err, redis.Nil) {
		return 0, errorpkg.ErrUTXONotFound
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *repo) SetBalance(ctx context.Context, userId string, amount float32) error {
	return r.rdb.Set(ctx, utxoKeyPrefix+userId, amount, 0).Err()
}

func (r *repo) Clear(ctx context.Context) error {
	keys, err := r.rdb.Keys(ctx, utxoKeyPrefix+"*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.rdb.Del(ctx, keys...).Err()
	}
	return nil
}
