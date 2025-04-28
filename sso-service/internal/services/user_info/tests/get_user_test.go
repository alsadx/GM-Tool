package tests

import (
	"context"
	"errors"
	"sso/internal/domain/models"
	"sso/tests/setup"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser_Success(t *testing.T) {
	service, _, mockUserProvider := setup.TestUser(t)

	ctx := context.Background()

	user := &models.User{
		Id:        int64(1),
		Name:      "testuser",
		Email:     "test@example.com",
		PassHash:  []byte("$2a$10$vR/gn5MPG9g5JVZPlhj"),
		FullName:  "",
		IsAdmin:   false,
		AvatarUrl: "",
	}

	mockUserProvider.EXPECT().
		UserById(ctx, int64(1)).
		Return(user, nil)

	mockUserProvider.EXPECT().
		UserByEmail(ctx, "test@example.com").
		Return(user, nil)

	resp, err := service.GetUserById(ctx, user.Id)
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	assert.Equal(t, user, resp)

	resp, err = service.GetUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	assert.Equal(t, user, resp)
}

func TestGetUser_UserNotFound(t *testing.T) {
	service, _, mockUserProvider := setup.TestUser(t)

	ctx := context.Background()

	mockUserProvider.EXPECT().
		UserById(ctx, int64(1)).
		Return(&models.User{}, models.ErrUserNotFound)

	mockUserProvider.EXPECT().
		UserByEmail(ctx, "test@example.com").
		Return(&models.User{}, models.ErrUserNotFound)

	resp, err := service.GetUserById(ctx, int64(1))
	require.Error(t, err)
	require.Empty(t, resp)

	assert.Equal(t, models.ErrUserNotFound, errors.Unwrap(err))

	resp, err = service.GetUserByEmail(ctx, "test@example.com")
	require.Error(t, err)
	require.Empty(t, resp)

	assert.Equal(t, models.ErrUserNotFound, errors.Unwrap(err))
}
