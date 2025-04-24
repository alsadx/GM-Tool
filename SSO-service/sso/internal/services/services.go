package services

import (
	"context"
	"sso/internal/domain/models"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email, name string, passHash []byte) (int64, error)
	UpdateUser(ctx context.Context, user *models.User) (error)
	DeleteUser(ctx context.Context, userId int64) error
	SetSession(ctx context.Context, userId int64, session models.Session) error
	DeleteSession(ctx context.Context, userId int64) error
}

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UserById(ctx context.Context, userId int64) (models.User, error)
	UserByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
	HealthCheck(ctx context.Context) (error)
}
