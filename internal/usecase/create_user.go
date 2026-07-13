package usecase

import (
	"context"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/user"
)

type CreateUserUseCase struct {
	userService user.Service
}

func NewCreateUserUseCase(userService user.Service) *CreateUserUseCase {
	return &CreateUserUseCase{userService: userService}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input user.CreateUserInput) (*ent.User, error) {
	return uc.userService.Create(ctx, input)
}
