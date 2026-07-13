package usecase

import (
	"context"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/user"
	errorpkg "github.com/luique16/quitocoin/internal/error"
	"github.com/luique16/quitocoin/internal/provider"
)

type LoginUseCase struct {
	repo   user.Repository
	hasher provider.PasswordHasher
	jwt    provider.JWTProvider
}

func NewLoginUseCase(repo user.Repository, hasher provider.PasswordHasher, jwt provider.JWTProvider) *LoginUseCase {
	return &LoginUseCase{repo: repo, hasher: hasher, jwt: jwt}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	User  *ent.User
	Token string
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	if input.Email == "" {
		return nil, errorpkg.ErrEmailRequired
	}
	if input.Password == "" {
		return nil, errorpkg.ErrPasswordRequired
	}

	u, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := uc.hasher.Compare(input.Password, u.Password); err != nil {
		return nil, errorpkg.ErrInvalidCredentials
	}

	token, err := uc.jwt.GenerateToken(u.ID, u.PublicID)
	if err != nil {
		return nil, errorpkg.ErrInternal
	}

	return &LoginOutput{User: u, Token: token}, nil
}
