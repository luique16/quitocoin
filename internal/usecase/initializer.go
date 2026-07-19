package usecase

import (
	"context"
	"fmt"

	"github.com/luique16/quitocoin/internal/domain/block"
	"github.com/luique16/quitocoin/internal/domain/transaction"
	"github.com/luique16/quitocoin/internal/domain/utxo"
)

type InitializerUseCase struct {
	blockService block.Service
	memPool      transaction.MemPool
	utxoService  utxo.Service
}

func NewInitializerUseCase(blockService block.Service, memPool transaction.MemPool, utxoService utxo.Service) *InitializerUseCase {
	return &InitializerUseCase{
		blockService: blockService,
		memPool:      memPool,
		utxoService:  utxoService,
	}
}

func (uc *InitializerUseCase) Execute(ctx context.Context) error {
	uc.memPool.Clear()

	uc.utxoService.Clear(ctx)

	chainLength, err := uc.blockService.GetChainLength(ctx)
	if err != nil {
		return err
	}

	if chainLength == 0 {
		fmt.Println("Creating genesis block...")
		_, err := uc.blockService.CreateGenesisBlock(ctx)

		if err != nil {
			return err
		}
	}

	chainLength, err = uc.blockService.GetChainLength(ctx)

	if err != nil {
		return err
	}

	blocks, err := uc.blockService.GetLastBlocks(ctx, chainLength)

	if err != nil {
		return err
	}

	for _, block := range blocks {
		if block.PreviousHash == "" {
			continue
		}

		err = uc.utxoService.Credit(ctx, block.Miner, float32(block.Reward))

		if err != nil {
			return err
		}

		for _, tx := range block.Transactions {
			err = uc.utxoService.Debit(ctx, tx.From, tx.Amount+1)

			if err != nil {
				return err
			}

			err = uc.utxoService.Credit(ctx, tx.To, tx.Amount)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
