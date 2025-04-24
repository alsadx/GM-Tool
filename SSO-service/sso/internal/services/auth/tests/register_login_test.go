package tests

import (
	"context"
	"errors"
	"sso/internal/domain/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {
	service, mockUserSaver, _, mockHasher := setupTest(t)

	ctx := context.Background()
	username := "testuser"
	email := "test@example.com"
	password := "testpassword"
	hashedPassword := []byte("$2a$10$vR/gn5MPG9g5JVZPlhj")

	mockHasher.EXPECT().
		Hash(password).
		Return(hashedPassword, nil)

	mockUserSaver.EXPECT().
		SaveUser(ctx, email, username, hashedPassword).
		Return(int64(1), nil)

	userId, err := service.Register(ctx, email, password, username)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), userId)
}

func TestRegister_AlreadyExists(t *testing.T) {
	service, mockUserSaver, _, mockHasher := setupTest(t)

	ctx := context.Background()
	username := "testuser"
	email := "test@example.com"
	password := "testpassword"
	hashedPassword := []byte("$2a$10$vR/gn5MPG9g5JVZPlhj")

	mockHasher.EXPECT().
		Hash(password).
		Return(hashedPassword, nil)

	mockUserSaver.EXPECT().
		SaveUser(ctx, email, username, hashedPassword).
		Return(int64(0), models.ErrUserExists)

	userId, err := service.Register(ctx, email, password, username)
	assert.Error(t, err)
	assert.Equal(t, models.ErrUserExists, errors.Unwrap(err))
	assert.Equal(t, int64(0), userId)
}

func TestLogin_Success(t *testing.T) {
	service, mockUserSaver, mockUserProvider, mockHasher := setupTest(t)

	ctx := context.Background()
	username := "testuser"
	email := "test@example.com"
	password := "testpassword"
	hashedPassword := []byte("$2a$10$vR/gn5MPG9g5JVZPlhj")

	mockUserProvider.EXPECT().
		UserByEmail(ctx, email).
		Return(models.User{
			Id:        1,
			Email:     email,
			PassHash:  hashedPassword,
			Name:      username,
			FullName:  "",
			IsAdmin:   false,
			AvatarUrl: "",
		}, nil)

	mockHasher.EXPECT().
		CheckPassword(password, hashedPassword).
		Return(nil)

	mockUserSaver.EXPECT().
		SetSession(ctx, int64(1), gomock.Any()).
		Return(nil)

	token, err := service.Login(ctx, email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
}
