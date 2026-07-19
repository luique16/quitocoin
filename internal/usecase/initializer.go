package usecase

import (
	"context"
	"fmt"

	"github.com/luique16/quitocoin/internal/domain/block"
	"github.com/luique16/quitocoin/internal/domain/transaction"
	"github.com/luique16/quitocoin/internal/domain/userblock"
	"github.com/luique16/quitocoin/internal/domain/utxo"
)

type InitializerUseCase struct {
	blockService     block.Service
	memPool          transaction.MemPool
	utxoService      utxo.Service
	userBlockService userblock.Service
}

func NewInitializerUseCase(blockService block.Service, memPool transaction.MemPool, utxoService utxo.Service, userBlockService userblock.Service) *InitializerUseCase {
	return &InitializerUseCase{
		blockService:     blockService,
		memPool:          memPool,
		utxoService:      utxoService,
		userBlockService: userBlockService,
	}
}

func (uc *InitializerUseCase) Execute(ctx context.Context) error {
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

	hasUTXO, err := uc.utxoService.HasData(ctx)
	if err != nil {
		return err
	}

	hasMempool := uc.memPool.Count() > 0

	hasUserBlocks, err := uc.userBlockService.HasData(ctx)
	if err != nil {
		return err
	}

	if hasUTXO && hasMempool && hasUserBlocks {
		return nil
	}

	fmt.Println("Loading blockchain into cache...")

	uc.utxoService.Clear(ctx)
	uc.memPool.Clear()
	uc.userBlockService.Clear(ctx)

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

		uc.userBlockService.AddBlock(ctx, block.Miner, block.Index, "miner")

		err = uc.utxoService.Credit(ctx, block.Miner, float32(block.Reward))

		if err != nil {
			return err
		}

		for _, tx := range block.Transactions {
			uc.userBlockService.AddBlock(ctx, tx.From, block.Index, "sender")
			uc.userBlockService.AddBlock(ctx, tx.To, block.Index, "receiver")

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
