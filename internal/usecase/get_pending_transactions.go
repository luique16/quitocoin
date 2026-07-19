package usecase

import (
	"context"

	"github.com/luique16/quitocoin/internal/domain/transaction"
)

type GetPendingTransactionsUseCase struct {
	memPool              transaction.MemPool
	transactionsPerBlock int
}

func NewGetPendingTransactionsUseCase(memPool transaction.MemPool, transactionsPerBlock int) *GetPendingTransactionsUseCase {
	return &GetPendingTransactionsUseCase{
		memPool:              memPool,
		transactionsPerBlock: transactionsPerBlock,
	}
}

type PendingTransaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

type GetPendingTransactionsOutput struct {
	Transactions []PendingTransaction `json:"transactions"`
}

func (uc *GetPendingTransactionsUseCase) Execute(_ context.Context) (*GetPendingTransactionsOutput, error) {
	txs := uc.memPool.PullFirstTransactions(uc.transactionsPerBlock)

	result := make([]PendingTransaction, len(txs))
	for i, tx := range txs {
		result[i] = PendingTransaction{
			From:   tx.From,
			To:     tx.To,
			Amount: tx.Amount,
		}
	}

	return &GetPendingTransactionsOutput{Transactions: result}, nil
}
