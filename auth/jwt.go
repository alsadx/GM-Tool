package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(token, secret string) (int, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	if !parsedToken.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	userId := int(claims["uid"].(float64))

	return userId, nil
}

func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("empty authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return token, nil
}

func IsTokenExpired(tokenString string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("exp field not found in token claims")
	}

	if int64(exp) < time.Now().Unix() {
		return fmt.Errorf("token is expired")
	}

	return nil
}
