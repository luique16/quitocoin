package block

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/luique16/quitocoin/ent"
	entblock "github.com/luique16/quitocoin/ent/block"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type Repository interface {
	Create(ctx context.Context, b *ent.Block) (*ent.Block, error)
	GetByHash(ctx context.Context, hash string) (*ent.Block, error)
	GetByIndex(ctx context.Context, index int) (*ent.Block, error)
	GetLast(ctx context.Context) (*ent.Block, error)
	GetLastN(ctx context.Context, n int) ([]*ent.Block, error)
	GetBlocksDescending(ctx context.Context, limit, offset int) ([]*ent.Block, error)
	List(ctx context.Context) ([]*ent.Block, error)
	Count(ctx context.Context) (int, error)
}

type repo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *repo {
	return &repo{client: client}
}

func (r *repo) Create(ctx context.Context, b *ent.Block) (*ent.Block, error) {
	created, err := r.client.Block.Create().
		SetHash(b.Hash).
		SetIndex(b.Index).
		SetPreviousHash(b.PreviousHash).
		SetNonce(b.Nonce).
		SetMiner(b.Miner).
		SetReward(b.Reward).
		SetTransactions(b.Transactions).
		SetCreatedAt(b.CreatedAt).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errorpkg.ErrBlockExists
		}
		return nil, errorpkg.ErrInternal
	}
	return created, nil
}

func (r *repo) GetByHash(ctx context.Context, hash string) (*ent.Block, error) {
	b, err := r.client.Block.Query().Where(entblock.Hash(hash)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorpkg.ErrBlockNotFound
		}
		return nil, errorpkg.ErrInternal
	}
	return b, nil
}

func (r *repo) GetByIndex(ctx context.Context, index int) (*ent.Block, error) {
	b, err := r.client.Block.Query().Where(entblock.IndexEQ(index)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorpkg.ErrBlockNotFound
		}
		return nil, errorpkg.ErrInternal
	}
	return b, nil
}

func (r *repo) GetLast(ctx context.Context) (*ent.Block, error) {
	b, err := r.client.Block.Query().
		Order(entblock.ByIndex(sql.OrderDesc())).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorpkg.ErrBlockNotFound
		}
		return nil, errorpkg.ErrInternal
	}
	return b, nil
}

func (r *repo) GetLastN(ctx context.Context, n int) ([]*ent.Block, error) {
	all, err := r.client.Block.Query().
		Order(entblock.ByIndex(sql.OrderDesc())).
		Limit(n).
		All(ctx)
	if err != nil {
		return nil, errorpkg.ErrInternal
	}
	for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
		all[i], all[j] = all[j], all[i]
	}
	return all, nil
}

func (r *repo) GetBlocksDescending(ctx context.Context, limit, offset int) ([]*ent.Block, error) {
	all, err := r.client.Block.Query().
		Order(entblock.ByIndex(sql.OrderDesc())).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, errorpkg.ErrInternal
	}
	return all, nil
}

func (r *repo) List(ctx context.Context) ([]*ent.Block, error) {
	all, err := r.client.Block.Query().Order(entblock.ByIndex()).All(ctx)
	if err != nil {
		return nil, errorpkg.ErrInternal
	}
	return all, nil
}

func (r *repo) Count(ctx context.Context) (int, error) {
	count, err := r.client.Block.Query().Count(ctx)
	if err != nil {
		return 0, errorpkg.ErrInternal
	}
	return count, nil
}
