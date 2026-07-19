package usecase

import (
	"github.com/luique16/quitocoin/internal/domain/transaction"
)

type GetMyPendingTransactionsUseCase struct {
	memPool transaction.MemPool
}

func NewGetMyPendingTransactionsUseCase(memPool transaction.MemPool) *GetMyPendingTransactionsUseCase {
	return &GetMyPendingTransactionsUseCase{memPool: memPool}
}

type MyPendingTransaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
	Type   string  `json:"type"`
}

type GetMyPendingTransactionsOutput struct {
	Transactions []MyPendingTransaction `json:"transactions"`
}

type GetMyPendingTransactionsInput struct {
	PublicID string
	Limit    int
}

func (uc *GetMyPendingTransactionsUseCase) Execute(input GetMyPendingTransactionsInput) *GetMyPendingTransactionsOutput {
	total := uc.memPool.Count()
	if total == 0 {
		return &GetMyPendingTransactionsOutput{Transactions: []MyPendingTransaction{}}
	}

	txs := uc.memPool.PullFirstTransactions(total)

	result := make([]MyPendingTransaction, 0, len(txs))
	for _, tx := range txs {
		if tx.From != input.PublicID && tx.To != input.PublicID {
			continue
		}

		t := "sent"
		if tx.To == input.PublicID {
			t = "received"
		}

		result = append(result, MyPendingTransaction{
			From:   tx.From,
			To:     tx.To,
			Amount: tx.Amount,
			Type:   t,
		})

		if input.Limit > 0 && len(result) >= input.Limit {
			break
		}
	}

	return &GetMyPendingTransactionsOutput{Transactions: result}
}
