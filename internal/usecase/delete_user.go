package usecase

import (
	"context"

	"github.com/luique16/quitocoin/internal/domain/user"
)

type DeleteUserUseCase struct {
	userService user.Service
}

func NewDeleteUserUseCase(userService user.Service) *DeleteUserUseCase {
	return &DeleteUserUseCase{userService: userService}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id string) error {
	return uc.userService.Delete(ctx, id)
}
