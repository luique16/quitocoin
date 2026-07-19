package usecase

import (
	"context"

	"github.com/luique16/quitocoin/internal/domain/transaction"
	"github.com/luique16/quitocoin/internal/domain/user"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type CreateTransferUseCase struct {
	userService user.Service
	memPool     transaction.MemPool
}

func NewCreateTransferUseCase(userService user.Service, memPool transaction.MemPool) *CreateTransferUseCase {
	return &CreateTransferUseCase{userService: userService, memPool: memPool}
}

type CreateTransferInput struct {
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

type CreateTransferOutput struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

func (uc *CreateTransferUseCase) Execute(ctx context.Context, from string, input CreateTransferInput) (*CreateTransferOutput, error) {
	if input.To == "" {
		return nil, errorpkg.ErrInvalidID
	}
	if input.Amount <= 0 {
		return nil, errorpkg.ErrNegativeAmount
	}
	if from == "" {
		return nil, errorpkg.ErrInvalidID
	}
	if from == input.To {
		return nil, errorpkg.ErrSelfTransfer
	}

	if _, err := uc.userService.GetByPublicID(ctx, input.To); err != nil {
		return nil, err
	}

	netAmount := input.Amount - 1
	tx := uc.memPool.CreateTransaction(from, input.To, netAmount)
	uc.memPool.PushTransaction(tx)

	return &CreateTransferOutput{
		From:   tx.From,
		To:     tx.To,
		Amount: tx.Amount,
	}, nil
}
