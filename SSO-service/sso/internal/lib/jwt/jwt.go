package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sso/internal/domain/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// type TokenManager interface {
// 	NewJWT(userId string, ttl time.Duration) (string, error)
// 	NewRefreshToken() (string, error)
// 	ParseJWT(token string) (string, error)
// }

// type Manager struct {
// 	// signKey string
// }

// func NewManager(signKey string) (*Manager, error) {
// 	// if signKey == "" {
// 	// 	return nil, fmt.Errorf("signing key should not be empty")
// 	// }
// 	return &Manager{
// 		// signKey: signKey,
// 	}, nil
// }

type TokenManager struct {
}

type ParsedJWT struct {
	UserId int64
	Email  string
	ExpiresAt time.Time
	AppId int
}
func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

func (m *TokenManager) NewJWT(user models.User, app models.App, ttl time.Duration) (string, error) {
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, ClaimsJWT{
	// 	UserId:   userId,
	// 	ExpiresAt: time.Now().Add(ttl).Unix(),
	// })

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = app.Id

	return token.SignedString([]byte(app.SigningKey))
}

func (m *TokenManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %s", err)
	}

	return hex.EncodeToString(b), nil
}

func (m *TokenManager) ParseJWT(token string, signKey string) (*ParsedJWT, error) {
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %s", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	var parsedJWT ParsedJWT
	parsedJWT.UserId = int64(claims["uid"].(float64))
	parsedJWT.Email = claims["email"].(string)
	parsedJWT.ExpiresAt = time.Unix(int64(claims["exp"].(float64)), 0)
	parsedJWT.AppId = int(claims["app_id"].(float64))

	// userId, ok := claims["sub"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("invalid token")
	// }

	return &parsedJWT, nil
}