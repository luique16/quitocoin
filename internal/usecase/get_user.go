package usecase

import (
	"context"
	"errors"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/user"
	"github.com/luique16/quitocoin/internal/domain/utxo"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type GetUserUseCase struct {
	userService user.Service
	utxoService utxo.Service
}

func NewGetUserUseCase(userService user.Service, utxoService utxo.Service) *GetUserUseCase {
	return &GetUserUseCase{userService: userService, utxoService: utxoService}
}

type GetUserOutput struct {
	User    *ent.User
	Balance float32
}

func (uc *GetUserUseCase) Execute(ctx context.Context, userID, publicID string) (*GetUserOutput, error) {
	u, err := uc.userService.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	balance, err := uc.utxoService.GetBalance(ctx, publicID)
	if err != nil && !errors.Is(err, errorpkg.ErrUTXONotFound) {
		return nil, err
	}

	return &GetUserOutput{User: u, Balance: balance}, nil
}
