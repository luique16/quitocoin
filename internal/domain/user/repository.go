package user

import (
	"context"

	"github.com/luique16/quitocoin/ent"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type Repository interface {
	Create(ctx context.Context, u *ent.User) (*ent.User, error)
	Get(ctx context.Context, id string) (*ent.User, error)
	List(ctx context.Context) ([]*ent.User, error)
	Update(ctx context.Context, u *ent.User) (*ent.User, error)
	Delete(ctx context.Context, id string) error
}

type repo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	return &repo{client: client}
}

func (r *repo) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	created, err := r.client.User.Create().
		SetID(u.ID).
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetPublicID(u.PublicID).
		SetCreatedAt(u.CreatedAt).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errorpkg.ErrEmailExists
		}
		return nil, errorpkg.ErrInternal
	}
	return created, nil
}

func (r *repo) Get(ctx context.Context, id string) (*ent.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorpkg.ErrUserNotFound
		}
		return nil, errorpkg.ErrInternal
	}
	return u, nil
}

func (r *repo) List(ctx context.Context) ([]*ent.User, error) {
	all, err := r.client.User.Query().All(ctx)
	if err != nil {
		return nil, errorpkg.ErrInternal
	}
	return all, nil
}

func (r *repo) Update(ctx context.Context, u *ent.User) (*ent.User, error) {
	updated, err := r.client.User.UpdateOneID(u.ID).
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetPublicID(u.PublicID).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorpkg.ErrUserNotFound
		}
		if ent.IsConstraintError(err) {
			return nil, errorpkg.ErrEmailExists
		}
		return nil, errorpkg.ErrInternal
	}
	return updated, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errorpkg.ErrUserNotFound
		}
		return errorpkg.ErrInternal
	}
	return nil
}
