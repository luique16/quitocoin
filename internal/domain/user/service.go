package user

import (
	"context"
	"strings"
	"time"

	"github.com/luique16/quitocoin/ent"
	errorpkg "github.com/luique16/quitocoin/internal/error"
	"github.com/luique16/quitocoin/internal/provider"
)

type Service interface {
	Create(ctx context.Context, input CreateUserInput) (*ent.User, error)
	Get(ctx context.Context, id string) (*ent.User, error)
	List(ctx context.Context) ([]*ent.User, error)
	Update(ctx context.Context, id string, input UpdateUserInput) (*ent.User, error)
	Delete(ctx context.Context, id string) error
}

func NewService(repo Repository, hasher provider.PasswordHasher, idGen provider.IDGenerator) Service {
	return &service{
		repo:   repo,
		hasher: hasher,
		idGen:  idGen,
	}
}

type service struct {
	repo   Repository
	hasher provider.PasswordHasher
	idGen  provider.IDGenerator
}

func (s *service) Create(ctx context.Context, input CreateUserInput) (*ent.User, error) {
	if err := validateCreate(input); err != nil {
		return nil, err
	}

	id := s.idGen.Generate()
	public_id := s.idGen.GeneratePublic()
	hashed, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, errorpkg.ErrInternal
	}

	u := &ent.User{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashed,
		PublicID:  public_id,
		CreatedAt: time.Now().UTC(),
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *service) Get(ctx context.Context, id string) (*ent.User, error) {
	if id == "" {
		return nil, errorpkg.ErrInvalidID
	}

	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) List(ctx context.Context) ([]*ent.User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateUserInput) (*ent.User, error) {
	if id == "" {
		return nil, errorpkg.ErrInvalidID
	}

	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	merged := &ent.User{
		ID:        existing.ID,
		Name:      existing.Name,
		Email:     existing.Email,
		Password:  existing.Password,
		PublicID:  existing.PublicID,
		CreatedAt: existing.CreatedAt,
	}

	if input.Name != nil {
		merged.Name = *input.Name
	}
	if input.Email != nil {
		merged.Email = *input.Email
	}
	if input.Password != nil {
		hashed, err := s.hasher.Hash(*input.Password)
		if err != nil {
			return nil, errorpkg.ErrInternal
		}
		merged.Password = hashed
	}

	updated, err := s.repo.Update(ctx, merged)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errorpkg.ErrInvalidID
	}

	if _, err := s.repo.Get(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// -- helpers -------------------------------------------------------------

func validateCreate(input CreateUserInput) error {
	if input.Name == "" {
		return errorpkg.ErrNameRequired
	}

	if input.Email == "" {
		return errorpkg.ErrEmailRequired
	}

	if !isValidEmail(input.Email) {
		return errorpkg.ErrInvalidEmail
	}

	if input.Password == "" {
		return errorpkg.ErrPasswordRequired
	}

	return nil
}

func isValidEmail(email string) bool {
	if strings.Count(email, "@") != 1 {
		return false
	}

	idx := strings.IndexByte(email, '@')

	if idx <= 0 || idx >= len(email)-1 {
		return false
	}

	if strings.ContainsAny(email, " ") {
		return false
	}

	return true
}
