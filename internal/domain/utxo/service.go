package utxo

import (
	"context"

	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type Service interface {
	SetBalance(ctx context.Context, userId string, amount float32) error
	GetBalance(ctx context.Context, userId string) (float32, error)
	Credit(ctx context.Context, userId string, amount float32) error
	Debit(ctx context.Context, userId string, amount float32) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SetBalance(ctx context.Context, userId string, amount float32) error {
	return s.repo.SetBalance(ctx, userId, amount)
}

func (s *service) GetBalance(ctx context.Context, userId string) (float32, error) {
	return s.repo.GetBalance(ctx, userId)
}

func (s *service) Credit(ctx context.Context, userId string, amount float32) error {
	if amount < 0 {
		return errorpkg.ErrNegativeAmount
	}
	if amount == 0 {
		return nil
	}

	balance, err := s.repo.GetBalance(ctx, userId)
	if err != nil {
		return err
	}

	return s.repo.SetBalance(ctx, userId, balance+amount)
}

func (s *service) Debit(ctx context.Context, userId string, amount float32) error {
	if amount < 0 {
		return errorpkg.ErrNegativeAmount
	}
	if amount == 0 {
		return nil
	}

	balance, err := s.repo.GetBalance(ctx, userId)
	if err != nil {
		return err
	}

	if balance < amount {
		return errorpkg.ErrInsufficientBalance
	}

	return s.repo.SetBalance(ctx, userId, balance-amount)
}
