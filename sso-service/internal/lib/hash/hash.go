package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (h *Hasher) CheckPassword(password string, hash []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("invalid password")
		}
		return fmt.Errorf("failed to compare passwords: %s", err)
	}

	return nil
}
