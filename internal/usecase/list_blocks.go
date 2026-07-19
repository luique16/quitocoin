package usecase

import (
	"context"

	"github.com/luique16/quitocoin/internal/domain/block"
)

type ListBlocksUseCase struct {
	blockService block.Service
}

func NewListBlocksUseCase(blockService block.Service) *ListBlocksUseCase {
	return &ListBlocksUseCase{blockService: blockService}
}

type BlockSummary struct {
	Index        int    `json:"index"`
	Hash         string `json:"hash"`
	Miner        string `json:"miner"`
	CreatedAt    string `json:"created_at"`
	TxCount      int    `json:"tx_count"`
}

type ListBlocksOutput struct {
	Blocks     []BlockSummary `json:"blocks"`
	TotalCount int            `json:"total_count"`
}

type ListBlocksInput struct {
	Limit  int
	Offset int
}

func (uc *ListBlocksUseCase) Execute(ctx context.Context, input ListBlocksInput) (*ListBlocksOutput, error) {
	total, err := uc.blockService.GetChainLength(ctx)
	if err != nil {
		return nil, err
	}

	blocks, err := uc.blockService.GetBlocksDescending(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	result := make([]BlockSummary, 0, len(blocks))
	for _, b := range blocks {
		result = append(result, BlockSummary{
			Index:     b.Index,
			Hash:      b.Hash,
			Miner:     b.Miner,
			CreatedAt: b.CreatedAt.Format("2006-01-02 15:04:05"),
			TxCount:   len(b.Transactions),
		})
	}

	return &ListBlocksOutput{
		Blocks:     result,
		TotalCount: total,
	}, nil
}
