package usecase

import (
	"context"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/user"
)

type ListUsersUseCase struct {
	userService user.Service
}

func NewListUsersUseCase(userService user.Service) *ListUsersUseCase {
	return &ListUsersUseCase{userService: userService}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context) ([]*ent.User, error) {
	return uc.userService.List(ctx)
}
