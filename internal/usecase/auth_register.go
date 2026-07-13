package usecase

import (
	"context"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/user"
	"github.com/luique16/quitocoin/internal/provider"
)

type RegisterUseCase struct {
	userService user.Service
	jwt         provider.JWTProvider
}

func NewRegisterUseCase(userService user.Service, jwt provider.JWTProvider) *RegisterUseCase {
	return &RegisterUseCase{userService: userService, jwt: jwt}
}

type RegisterOutput struct {
	User  *ent.User
	Token string
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input user.CreateUserInput) (*RegisterOutput, error) {
	u, err := uc.userService.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	token, err := uc.jwt.GenerateToken(u.ID, u.PublicID)
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{User: u, Token: token}, nil
}
