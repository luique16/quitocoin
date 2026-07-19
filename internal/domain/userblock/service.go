package userblock

import "context"

type Service interface {
	AddBlock(ctx context.Context, publicID string, blockIndex int, role string) error
	GetBlocks(ctx context.Context, publicID string, role string, limit int) ([]BlockRef, error)
	HasData(ctx context.Context) (bool, error)
	Clear(ctx context.Context) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddBlock(ctx context.Context, publicID string, blockIndex int, role string) error {
	return s.repo.AddBlock(ctx, publicID, blockIndex, role)
}

func (s *service) GetBlocks(ctx context.Context, publicID string, role string, limit int) ([]BlockRef, error) {
	return s.repo.GetBlocks(ctx, publicID, role, limit)
}

func (s *service) HasData(ctx context.Context) (bool, error) {
	return s.repo.HasData(ctx)
}

func (s *service) Clear(ctx context.Context) error {
	return s.repo.Clear(ctx)
}
