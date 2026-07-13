package provider

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) error
}

type passwordHasher struct {
	cost int
}

func NewPasswordHasher() PasswordHasher {
	return &passwordHasher{cost: bcrypt.DefaultCost}
}

func (h *passwordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *passwordHasher) Compare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
