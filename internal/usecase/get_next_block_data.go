package usecase

import (
	"context"
	"fmt"

	"github.com/luique16/quitocoin/internal/domain/block"
	"github.com/luique16/quitocoin/internal/domain/transaction"
)

type GetNextBlockDataUseCase struct {
	blockService         block.Service
	memPool              transaction.MemPool
	transactionsPerBlock int
}

func NewGetNextBlockDataUseCase(blockService block.Service, memPool transaction.MemPool, transactionsPerBlock int) *GetNextBlockDataUseCase {
	return &GetNextBlockDataUseCase{
		blockService:         blockService,
		memPool:              memPool,
		transactionsPerBlock: transactionsPerBlock,
	}
}

type NextBlockDataOutput struct {
	Text string `json:"text"`
}

func (uc *GetNextBlockDataUseCase) Execute(ctx context.Context, minerID string) (*NextBlockDataOutput, error) {
	last, err := uc.blockService.GetLastBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("get last block: %w", err)
	}

	nextIndex := last.Index + 1

	txs := uc.memPool.PullFirstTransactions(uc.transactionsPerBlock)
	reward := block.BaseReward + block.RewardPerTransaction*float32(len(txs))

	text := block.FormatBlockInput(nextIndex, minerID, reward, last.Hash, txs)

	return &NextBlockDataOutput{Text: text}, nil
}
