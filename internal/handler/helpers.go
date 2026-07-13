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
	CreatedAt time.Time `json:"created_at"`
}

func toUserResponse(u *ent.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		PublicID:  u.PublicID,
		CreatedAt: u.CreatedAt,
	}
}

func mapError(err error) int {
	switch {
		case errors.Is(err, errorpkg.ErrNameRequired),
			errors.Is(err, errorpkg.ErrEmailRequired),
			errors.Is(err, errorpkg.ErrInvalidEmail),
			errors.Is(err, errorpkg.ErrPasswordRequired),
			errors.Is(err, errorpkg.ErrInvalidID):
			return http.StatusBadRequest
		case errors.Is(err, errorpkg.ErrUserNotFound):
			return http.StatusNotFound
		case errors.Is(err, errorpkg.ErrEmailExists):
			return http.StatusConflict
		default:
			return http.StatusInternalServerError
	}
}
