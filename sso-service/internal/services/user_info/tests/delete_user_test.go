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

func TestDeleteUser_Success(t *testing.T) {
	service, mockUserSaver, _ := setup.TestUser(t)

	ctx := context.Background()

	userId := int64(1)

	mockUserSaver.EXPECT().
		DeleteSession(ctx, userId).
		Return(nil)

	mockUserSaver.EXPECT().
		DeleteUser(ctx, userId).
		Return(nil)

	err := service.DeleteUser(ctx, userId)
	require.NoError(t, err)
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	service, mockUserSaver, _ := setup.TestUser(t)

	ctx := context.Background()

	userId := int64(1)

	mockUserSaver.EXPECT().
		DeleteSession(ctx, userId).
		Return(models.ErrUserNotFound)

	err := service.DeleteUser(ctx, userId)
	require.Error(t, err)
	assert.Equal(t, models.ErrUserNotFound, errors.Unwrap(err))
}

func TestDeleteUser_UserNotFound2(t *testing.T) {
	service, mockUserSaver, _ := setup.TestUser(t)

	ctx := context.Background()

	userId := int64(1)

	mockUserSaver.EXPECT().
		DeleteSession(ctx, userId).
		Return(nil)

	mockUserSaver.EXPECT().
		DeleteUser(ctx, userId).
		Return(models.ErrUserNotFound)

	err := service.DeleteUser(ctx, userId)
	require.Error(t, err)
	assert.Equal(t, models.ErrUserNotFound, errors.Unwrap(err))
}
