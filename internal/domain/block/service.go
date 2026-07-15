package block

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/transaction"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type Service interface {
	CreateGenesisBlock(ctx context.Context) (*ent.Block, error)
	TryToMineBlock(ctx context.Context, miner string, nonce int64, reward float32) (*ent.Block, error)
	GetBlockByHash(ctx context.Context, hash string) (*ent.Block, error)
	GetBlockByIndex(ctx context.Context, index int) (*ent.Block, error)
	GetLastBlock(ctx context.Context) (*ent.Block, error)
	GetLastBlocks(ctx context.Context, n int) ([]*ent.Block, error)
	GetChainLength(ctx context.Context) (int, error)
	ValidateChain(ctx context.Context) (bool, error)
}

var DifficultyPrefix = strings.Repeat("0", Difficulty)

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
	var nonce int64
	now := time.Now().UTC()

	for {
		hash := CalculateHash(0, "system", 0.0, "", nonce, nil)
		if strings.HasPrefix(hash, DifficultyPrefix) {
			b := &ent.Block{
				Hash:         hash,
				Index:        0,
				PreviousHash: "",
				Nonce:        nonce,
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
		nonce++
	}
}

func (s *service) TryToMineBlock(ctx context.Context, miner string, nonce int64, reward float32) (*ent.Block, error) {
	last, err := s.repo.GetLast(ctx)
	if err != nil {
		return nil, err
	}

	newIndex := last.Index + 1
	hash := CalculateHash(newIndex, miner, reward, last.Hash, nonce, nil)

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
		CreatedAt:    time.Now().UTC(),
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
		expectedHash := CalculateHash(b.Index, b.Miner, float32(b.Reward), b.PreviousHash, b.Nonce, b.Transactions)
		if b.Hash != expectedHash {
			return false, nil
		}
		if i > 0 && b.PreviousHash != blocks[i-1].Hash {
			return false, nil
		}
	}

	return true, nil
}

func CalculateHash(index int, miner string, reward float32, previousHash string, nonce int64, transactions []transaction.Transaction) string {
	transactionsData := ""
	for i := range transactions {
		transactionsData += fmt.Sprintf("%s:%f:%s;", transactions[i].From, transactions[i].Amount, transactions[i].To)
	}

	data := fmt.Sprintf("%d:%s:%f:%s:%d:%s", index, miner, reward, previousHash, nonce, transactionsData)

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}
