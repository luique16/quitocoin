package provider

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	PublicID string `json:"public_id"`
	jwt.RegisteredClaims
}

type JWTProvider interface {
	GenerateToken(userID, publicID string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type jwtProvider struct {
	secret []byte
}

func NewJWTProvider(secret string) JWTProvider {
	return &jwtProvider{secret: []byte(secret)}
}

func (p *jwtProvider) GenerateToken(userID, publicID string) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		PublicID: publicID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *jwtProvider) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		return p.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
