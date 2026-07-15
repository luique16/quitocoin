package block

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/luique16/quitocoin/ent"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type Service interface {
	CreateGenesisBlock(ctx context.Context) (*ent.Block, error)
	TryToMineBlock(ctx context.Context, miner string, nonce int64) (*ent.Block, error)
	GetBlockByHash(ctx context.Context, hash string) (*ent.Block, error)
	GetBlockByIndex(ctx context.Context, index int) (*ent.Block, error)
	GetLastBlock(ctx context.Context) (*ent.Block, error)
	GetLastBlocks(ctx context.Context, n int) ([]*ent.Block, error)
	GetChainLength(ctx context.Context) (int, error)
	ValidateChain(ctx context.Context) (bool, error)
}

var DifficultyPrefix = "000"

type service struct {
	transactionsPerBlock int
	repo                 Repository
}

func NewService(transactionsPerBlock int, repo Repository) Service {
	return &service{
		transactionsPerBlock: transactionsPerBlock,
		repo:                 repo,
	}
}

func (s *service) CreateGenesisBlock(ctx context.Context) (*ent.Block, error) {
	now := time.Now().UTC()

	b := &ent.Block{
		Hash:         CalculateHash(0, now, "", 0, nil),
		Index:        0,
		PreviousHash: "",
		Nonce:        0,
		Miner:        "system",
		Reward:       0,
		Transactions: nil,
		CreatedAt:    now,
	}

	created, err := s.repo.Create(ctx, b)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *service) TryToMineBlock(ctx context.Context, miner string, nonce int64) (*ent.Block, error) {
	last, err := s.repo.GetLast(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	newIndex := last.Index + 1
	hash := CalculateHash(newIndex, now, last.Hash, nonce, nil)

	if !strings.HasPrefix(hash, DifficultyPrefix) {
		return nil, errorpkg.ErrInvalidNonce
	}

	b := &ent.Block{
		Hash:         hash,
		Index:        newIndex,
		PreviousHash: last.Hash,
		Nonce:        nonce,
		Miner:        miner,
		Reward:       1.0,
		Transactions: nil,
		CreatedAt:    now,
	}

	created, err := s.repo.Create(ctx, b)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *service) GetBlockByHash(ctx context.Context, hash string) (*ent.Block, error) {
	if hash == "" {
		return nil, errorpkg.ErrInvalidBlockHash
	}
	return s.repo.GetByHash(ctx, hash)
}

func (s *service) GetBlockByIndex(ctx context.Context, index int) (*ent.Block, error) {
	if index < 0 {
		return nil, errorpkg.ErrInvalidBlockIndex
	}
	return s.repo.GetByIndex(ctx, index)
}

func (s *service) GetLastBlock(ctx context.Context) (*ent.Block, error) {
	return s.repo.GetLast(ctx)
}

func (s *service) GetLastBlocks(ctx context.Context, n int) ([]*ent.Block, error) {
	if n <= 0 {
		return []*ent.Block{}, nil
	}
	return s.repo.GetLastN(ctx, n)
}

func (s *service) GetChainLength(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s *service) ValidateChain(ctx context.Context) (bool, error) {
	blocks, err := s.repo.List(ctx)
	if err != nil {
		return false, err
	}

	if len(blocks) == 0 {
		return true, nil
	}

	for i, b := range blocks {
		expectedHash := CalculateHash(b.Index, b.CreatedAt, b.PreviousHash, b.Nonce, b.Transactions)
		if b.Hash != expectedHash {
			return false, nil
		}
		if i > 0 && b.PreviousHash != blocks[i-1].Hash {
			return false, nil
		}
	}

	return true, nil
}

func CalculateHash(index int, timestamp time.Time, previousHash string, nonce int64, transactions any) string {
	data := fmt.Sprintf("%d:%d:%s:%d:%v", index, timestamp.UnixNano(), previousHash, nonce, transactions)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
