package usecase

import (
	"context"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/user"
)

type UpdateUserUseCase struct {
	userService user.Service
}

func NewUpdateUserUseCase(userService user.Service) *UpdateUserUseCase {
	return &UpdateUserUseCase{userService: userService}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, id string, input user.UpdateUserInput) (*ent.User, error) {
	return uc.userService.Update(ctx, id, input)
}
