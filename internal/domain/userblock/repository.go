package userblock

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const userBlocksPrefix = "user_blocks:"

type Repository interface {
	AddBlock(ctx context.Context, publicID string, blockIndex int) error
	GetBlocks(ctx context.Context, publicID string) ([]int, error)
	Clear(ctx context.Context) error
}

type repo struct {
	rdb redis.Cmdable
}

func NewRepository(rdb redis.Cmdable) Repository {
	return &repo{rdb: rdb}
}

func (r *repo) AddBlock(ctx context.Context, publicID string, blockIndex int) error {
	return r.rdb.SAdd(ctx, userBlocksPrefix+publicID, blockIndex).Err()
}

func (r *repo) GetBlocks(ctx context.Context, publicID string) ([]int, error) {
	members, err := r.rdb.SMembers(ctx, userBlocksPrefix+publicID).Result()
	if err != nil {
		return nil, err
	}

	blocks := make([]int, 0, len(members))
	for _, m := range members {
		v, err := strconv.Atoi(m)
		if err != nil {
			continue
		}
		blocks = append(blocks, v)
	}
	return blocks, nil
}

func (r *repo) Clear(ctx context.Context) error {
	keys, err := r.rdb.Keys(ctx, userBlocksPrefix+"*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.rdb.Del(ctx, keys...).Err()
	}
	return nil
}
