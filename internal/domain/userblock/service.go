package userblock

import "context"

type Service interface {
	AddBlock(ctx context.Context, publicID string, blockIndex int) error
	GetBlocks(ctx context.Context, publicID string) ([]int, error)
	Clear(ctx context.Context) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddBlock(ctx context.Context, publicID string, blockIndex int) error {
	return s.repo.AddBlock(ctx, publicID, blockIndex)
}

func (s *service) GetBlocks(ctx context.Context, publicID string) ([]int, error) {
	return s.repo.GetBlocks(ctx, publicID)
}

func (s *service) Clear(ctx context.Context) error {
	return s.repo.Clear(ctx)
}
