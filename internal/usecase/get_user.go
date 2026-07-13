package usecase

import (
	"context"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/user"
)

type GetUserUseCase struct {
	userService user.Service
}

func NewGetUserUseCase(userService user.Service) *GetUserUseCase {
	return &GetUserUseCase{userService: userService}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*ent.User, error) {
	return uc.userService.Get(ctx, id)
}
