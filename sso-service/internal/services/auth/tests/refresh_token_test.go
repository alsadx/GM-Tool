package tests

import (
	"context"
	"errors"
	"sso/internal/domain/models"
	"sso/tests/setup"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken_Success(t *testing.T) {
	service, mockUserSaver, mockUserProvider, _, mockTokenManager := setup.TestAuth(t)
	ctx := context.Background()

	refreshToken := "refresh_token"

	user := models.User{
		Id:        1,
		Email:     "email",
		PassHash:  []byte("pass_hash"),
		Name:      "name",
		FullName:  "full_name",
		IsAdmin:   false,
		AvatarUrl: "avatar_url",
	}

	mockUserProvider.EXPECT().
		UserByRefreshToken(ctx, refreshToken).
		Return(user, nil)

	mockTokenManager.EXPECT().
		NewJWT(user, "secret", gomock.Any()).
		Return("access_token", nil)

	mockTokenManager.EXPECT().
		NewRefreshToken().
		Return("new_refresh_token", nil)

	mockUserSaver.EXPECT().
		SetSession(ctx, user.Id, gomock.Any()).
		Return(nil)

	tokens, err := service.RefreshToken(ctx, refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens)
	assert.Equal(t, "access_token", tokens.AccessToken)
	assert.Equal(t, "new_refresh_token", tokens.RefreshToken)
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	service, _, mockUserProvider, _, _ := setup.TestAuth(t)
	ctx := context.Background()

	refreshToken := "refresh_token"

	mockUserProvider.EXPECT().
		UserByRefreshToken(ctx, refreshToken).
		Return(models.User{}, models.ErrUserNotFound)

	tokens, err := service.RefreshToken(ctx, refreshToken)
	assert.Error(t, err)
	assert.Empty(t, tokens)
	assert.Equal(t, errors.Unwrap(err), models.ErrInvalidRefreshToken)
}
