package userblock

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

const userBlocksPrefix = "user_blocks:"

type BlockRef struct {
	Index int
	Role  string
}

type Repository interface {
	AddBlock(ctx context.Context, publicID string, blockIndex int, role string) error
	GetBlocks(ctx context.Context, publicID string, role string, limit int) ([]BlockRef, error)
	HasData(ctx context.Context) (bool, error)
	Clear(ctx context.Context) error
}

type repo struct {
	rdb redis.Cmdable
}

func NewRepository(rdb redis.Cmdable) Repository {
	return &repo{rdb: rdb}
}

func (r *repo) AddBlock(ctx context.Context, publicID string, blockIndex int, role string) error {
	member := fmt.Sprintf("%s:%d", role, blockIndex)
	return r.rdb.ZAdd(ctx, userBlocksPrefix+publicID, redis.Z{
		Score:  float64(blockIndex),
		Member: member,
	}).Err()
}

func (r *repo) GetBlocks(ctx context.Context, publicID string, role string, limit int) ([]BlockRef, error) {
	members, err := r.rdb.ZRevRange(ctx, userBlocksPrefix+publicID, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	result := make([]BlockRef, 0, len(members))
	for _, m := range members {
		parts := strings.SplitN(m, ":", 2)
		if len(parts) != 2 {
			continue
		}
		memberRole := parts[0]
		index, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		if role != "" && memberRole != role {
			continue
		}

		result = append(result, BlockRef{Index: index, Role: memberRole})
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func (r *repo) HasData(ctx context.Context) (bool, error) {
	keys, err := r.rdb.Keys(ctx, userBlocksPrefix+"*").Result()
	if err != nil {
		return false, err
	}
	if len(keys) == 0 {
		return false, nil
	}

	members, err := r.rdb.ZRevRange(ctx, keys[0], 0, -1).Result()
	if err != nil {
		return false, nil
	}

	for _, m := range members {
		parts := strings.SplitN(m, ":", 2)
		if len(parts) == 2 {
			if _, err := strconv.Atoi(parts[1]); err == nil {
				return true, nil
			}
		}
	}

	return false, nil
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
