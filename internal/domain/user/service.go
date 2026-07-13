package user

import (
	"context"

	"github.com/luique16/quitocoin/ent"
)

type Service interface {
	Create(ctx context.Context, input CreateUserInput) (*ent.User, error)
	Get(ctx context.Context, id string) (*ent.User, error)
	List(ctx context.Context) ([]*ent.User, error)
	Update(ctx context.Context, id string, input UpdateUserInput) (*ent.User, error)
	Delete(ctx context.Context, id string) error
}
