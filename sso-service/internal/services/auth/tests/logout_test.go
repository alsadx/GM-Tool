package tests

import (
	"context"
	"errors"
	"sso/internal/domain/models"
	"sso/tests/setup"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogout_Success(t *testing.T) {
	service, mockUserSaver, _, _, _ := setup.TestAuth(t)

	ctx := context.Background()
	userId := int64(1)

	mockUserSaver.EXPECT().
		DeleteSession(ctx, userId).
		Return(nil)

	err := service.Logout(ctx, userId)
	assert.NoError(t, err)
}

func TestLogout_UserNotFound(t *testing.T) {
	service, mockUserSaver, _, _, _ := setup.TestAuth(t)

	ctx := context.Background()
	userId := int64(1)

	mockUserSaver.EXPECT().
		DeleteSession(ctx, userId).
		Return(models.ErrUserNotFound)

	err := service.Logout(ctx, userId)
	assert.Error(t, err)
	assert.Equal(t, models.ErrUserNotFound, errors.Unwrap(err))
}
