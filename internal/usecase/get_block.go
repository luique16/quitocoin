package usecase

import (
	"context"
	"errors"

	"github.com/luique16/quitocoin/internal/domain/block"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type GetBlockUseCase struct {
	blockService block.Service
}

func NewGetBlockUseCase(blockService block.Service) *GetBlockUseCase {
	return &GetBlockUseCase{blockService: blockService}
}

type BlockTx struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

type BlockDetail struct {
	Index        int       `json:"index"`
	Hash         string    `json:"hash"`
	PreviousHash string    `json:"previous_hash"`
	Nonce        int64     `json:"nonce"`
	Miner        string    `json:"miner"`
	Reward       float64   `json:"reward"`
	CreatedAt    string    `json:"created_at"`
	TxCount      int       `json:"tx_count"`
	Transactions []BlockTx `json:"transactions"`
}

type GetBlockOutput struct {
	Block BlockDetail `json:"block"`
}

type GetBlockInput struct {
	Index int
}

func (uc *GetBlockUseCase) Execute(ctx context.Context, input GetBlockInput) (*GetBlockOutput, error) {
	b, err := uc.blockService.GetBlockByIndex(ctx, input.Index)
	if err != nil {
		if errors.Is(err, errorpkg.ErrBlockNotFound) {
			return nil, errorpkg.ErrBlockNotFound
		}
		return nil, err
	}

	txs := make([]BlockTx, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		txs = append(txs, BlockTx{
			From:   tx.From,
			To:     tx.To,
			Amount: tx.Amount,
		})
	}

	return &GetBlockOutput{
		Block: BlockDetail{
			Index:        b.Index,
			Hash:         b.Hash,
			PreviousHash: b.PreviousHash,
			Nonce:        b.Nonce,
			Miner:        b.Miner,
			Reward:       b.Reward,
			CreatedAt:    b.CreatedAt.Format("2006-01-02 15:04:05"),
			TxCount:      len(b.Transactions),
			Transactions: txs,
		},
	}, nil
}
