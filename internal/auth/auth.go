package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	NewRefreshToken() (string, error)
	ParseJWT(token string) (string, error)
}

type Manager struct {
	signKey string
}

func NewManager(signKey string) (*Manager, error) {
	if signKey == "" {
		return nil, fmt.Errorf("signing key should not be empty")
	}
	return &Manager{
		signKey: signKey,
	}, nil
}

func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	})

	return token.SignedString([]byte(m.signKey))
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %s", err)
	}

	return hex.EncodeToString(b), nil
}

func (m *Manager) ParseJWT(token string) (string, error) {
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(m.signKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %s", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return "", fmt.Errorf("invalid token")
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}

	return userId, nil
}