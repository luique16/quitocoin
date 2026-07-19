package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/luique16/quitocoin/ent"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

type userResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PublicID  string    `json:"public_id"`
	Balance   float32   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserResponse(u *ent.User, balance float32) userResponse {
	return userResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		PublicID:  u.PublicID,
		Balance:   balance,
		CreatedAt: u.CreatedAt,
	}
}

func mapError(err error) int {
	switch {
	case errors.Is(err, errorpkg.ErrNameRequired),
		errors.Is(err, errorpkg.ErrEmailRequired),
		errors.Is(err, errorpkg.ErrInvalidEmail),
		errors.Is(err, errorpkg.ErrPasswordRequired),
		errors.Is(err, errorpkg.ErrWeakPassword),
		errors.Is(err, errorpkg.ErrInvalidID),
		errors.Is(err, errorpkg.ErrNegativeAmount),
		errors.Is(err, errorpkg.ErrInvalidBlockHash),
		errors.Is(err, errorpkg.ErrInvalidNonce),
		errors.Is(err, errorpkg.ErrInvalidBlockIndex),
		errors.Is(err, errorpkg.ErrSelfTransfer):
		return http.StatusBadRequest
	case errors.Is(err, errorpkg.ErrInvalidCredentials),
		errors.Is(err, errorpkg.ErrIncorrectPassword):
		return http.StatusUnauthorized
	case errors.Is(err, errorpkg.ErrUserNotFound),
		errors.Is(err, errorpkg.ErrBlockNotFound),
		errors.Is(err, errorpkg.ErrUTXONotFound):
		return http.StatusNotFound
	case errors.Is(err, errorpkg.ErrEmailExists),
		errors.Is(err, errorpkg.ErrBlockExists):
		return http.StatusConflict
	case errors.Is(err, errorpkg.ErrInsufficientBalance):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
