package utxo

import (
	"context"
	"errors"

	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type Service interface {
	SetBalance(ctx context.Context, userId string, amount float32) error
	GetBalance(ctx context.Context, userId string) (float32, error)
	GetAll(ctx context.Context) ([]Entry, error)
	HasData(ctx context.Context) (bool, error)
	Credit(ctx context.Context, userId string, amount float32) error
	Debit(ctx context.Context, userId string, amount float32) error
	Clear(ctx context.Context) error
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
	if err != nil && !errors.Is(err, errorpkg.ErrUTXONotFound) {
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
	if err != nil && !errors.Is(err, errorpkg.ErrUTXONotFound) {
		return err
	}

	if balance < amount {
		return errorpkg.ErrInsufficientBalance
	}

	return s.repo.SetBalance(ctx, userId, balance-amount)
}

func (s *service) Clear(ctx context.Context) error {
	return s.repo.Clear(ctx)
}

func (s *service) GetAll(ctx context.Context) ([]Entry, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) HasData(ctx context.Context) (bool, error) {
	return s.repo.HasData(ctx)
}
