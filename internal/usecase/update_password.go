package usecase

import (
	"context"

	"github.com/luique16/quitocoin/internal/domain/user"
)

type UpdatePasswordUseCase struct {
	userService user.Service
}

func NewUpdatePasswordUseCase(userService user.Service) *UpdatePasswordUseCase {
	return &UpdatePasswordUseCase{userService: userService}
}

func (uc *UpdatePasswordUseCase) Execute(ctx context.Context, id string, input user.UpdatePasswordInput) error {
	return uc.userService.UpdatePassword(ctx, id, input)
}
