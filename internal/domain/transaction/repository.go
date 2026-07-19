package transaction

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

const mempoolKey = "mempool"

type Repository interface {
	Push(ctx context.Context, t Transaction) error
	PullFirst(ctx context.Context, n int) ([]Transaction, error)
	DeleteFirst(ctx context.Context, n int) error
	Count(ctx context.Context) (int, error)
	Clear(ctx context.Context) error
}

type repo struct {
	rdb redis.Cmdable
}

func NewRepository(rdb redis.Cmdable) Repository {
	return &repo{rdb: rdb}
}

func (r *repo) Push(ctx context.Context, t Transaction) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return r.rdb.RPush(ctx, mempoolKey, data).Err()
}

func (r *repo) PullFirst(ctx context.Context, n int) ([]Transaction, error) {
	if n <= 0 {
		return []Transaction{}, nil
	}

	data, err := r.rdb.LRange(ctx, mempoolKey, 0, int64(n-1)).Result()
	if err != nil {
		return nil, err
	}

	result := make([]Transaction, 0, len(data))
	for _, item := range data {
		var tx Transaction
		if err := json.Unmarshal([]byte(item), &tx); err != nil {
			continue
		}
		result = append(result, tx)
	}
	return result, nil
}

func (r *repo) DeleteFirst(ctx context.Context, n int) error {
	if n <= 0 {
		return nil
	}
	return r.rdb.LTrim(ctx, mempoolKey, int64(n), -1).Err()
}

func (r *repo) Count(ctx context.Context) (int, error) {
	count, err := r.rdb.LLen(ctx, mempoolKey).Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *repo) Clear(ctx context.Context) error {
	return r.rdb.Del(ctx, mempoolKey).Err()
}
