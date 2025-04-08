package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return string(bytes)
}
