package usecase

import (
	"context"
	"fmt"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/block"
	"github.com/luique16/quitocoin/internal/domain/transaction"
	"github.com/luique16/quitocoin/internal/domain/userblock"
	"github.com/luique16/quitocoin/internal/domain/utxo"
)

type MineBlockUseCase struct {
	blockService         block.Service
	utxoService          utxo.Service
	memPool              transaction.MemPool
	userBlockService     userblock.Service
	transactionsPerBlock int
}

func NewMineBlockUseCase(blockService block.Service, utxoService utxo.Service, memPool transaction.MemPool, userBlockService userblock.Service, transactionsPerBlock int) *MineBlockUseCase {
	return &MineBlockUseCase{
		blockService:         blockService,
		utxoService:          utxoService,
		memPool:              memPool,
		userBlockService:     userBlockService,
		transactionsPerBlock: transactionsPerBlock,
	}
}

type MineBlockInput struct {
	Nonce int64 `json:"nonce"`
}

func (uc *MineBlockUseCase) Execute(ctx context.Context, minerID string, input MineBlockInput) (*ent.Block, error) {
	txs := uc.memPool.PullFirstTransactions(uc.transactionsPerBlock)
	reward := block.BaseReward + block.RewardPerTransaction*float32(len(txs))

	b, err := uc.blockService.TryToMineBlock(ctx, minerID, input.Nonce, reward, txs)
	if err != nil {
		return nil, err
	}

	uc.memPool.DeleteFirstTransactions(len(txs))

	if err := uc.utxoService.Credit(ctx, minerID, reward); err != nil {
		return nil, fmt.Errorf("credit miner: %w", err)
	}

	uc.userBlockService.AddBlock(ctx, minerID, b.Index, "miner")

	for _, tx := range txs {
		if err := uc.utxoService.Debit(ctx, tx.From, tx.Amount+1); err != nil {
			return nil, fmt.Errorf("debit sender %s: %w", tx.From, err)
		}
		if err := uc.utxoService.Credit(ctx, tx.To, tx.Amount); err != nil {
			return nil, fmt.Errorf("credit receiver %s: %w", tx.To, err)
		}

		uc.userBlockService.AddBlock(ctx, tx.From, b.Index, "sender")
		uc.userBlockService.AddBlock(ctx, tx.To, b.Index, "receiver")
	}

	return b, nil
}
