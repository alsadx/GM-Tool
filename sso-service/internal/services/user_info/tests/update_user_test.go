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

func TestUpdateUser_Success(t *testing.T) {
	service, mockUserSaver, mockUserProvider := setup.TestUser(t)

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

	updatedUser := &models.User{
		Id:        int64(1),
		Name:      "testuser2",
		Email:     "test@example.com",
		PassHash:  []byte("$2a$10$vR/gn5MPG9g5JVZPlhj"),
		FullName:  "Test User 2",
		IsAdmin:   false,
		AvatarUrl: "path/to/avatar.jpg",
	}

	updates := map[string]string{
		"name":      "testuser2",
		"full_name": "Test User 2",
		"avatar_url":    "path/to/avatar.jpg",
	}

	mockUserProvider.EXPECT().
		UserById(ctx, int64(1)).
		Return(user, nil)

	mockUserSaver.EXPECT().
		UpdateUser(ctx, user).
		Return(nil)

	resp, err := service.UpdateUser(ctx, user.Id, updates)
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	assert.Equal(t, updatedUser, resp)
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	service, _, mockUserProvider := setup.TestUser(t)

	ctx := context.Background()

	updates := map[string]string{
		"name":      "testuser2",
		"full_name": "Test User 2",
		"avatar_url":    "path/to/avatar.jpg",
	}

	mockUserProvider.EXPECT().
		UserById(ctx, int64(1)).
		Return(&models.User{}, models.ErrUserNotFound)

	resp, err := service.UpdateUser(ctx, int64(1), updates)
	require.Error(t, err)
	require.Empty(t, resp)

	assert.Equal(t, models.ErrUserNotFound, errors.Unwrap(err))
}

func TestUpdateUser_InvalidArguments(t *testing.T) {
	service, mockUserSaver, mockUserProvider := setup.TestUser(t)

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

	updates := map[string]string{
		"name":      "testuser2",
		"full_name": "Test User 2",
		"avatar":    "path/to/avatar.jpg",
	}

	mockUserProvider.EXPECT().
		UserById(ctx, int64(1)).
		Return(user, nil)

	mockUserSaver.EXPECT().
		UpdateUser(ctx, user).
		Return(models.ErrInvalidArgument)

	resp, err := service.UpdateUser(ctx, user.Id, updates)
	require.Error(t, err)
	require.Empty(t, resp)

	assert.Equal(t, models.ErrInvalidArgument, errors.Unwrap(err))
}

func TestUpdateUser_EmptyUpdates(t *testing.T) {
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

	updates := map[string]string{}

	mockUserProvider.EXPECT().
		UserById(ctx, int64(1)).
		Return(user, nil)

	resp, err := service.UpdateUser(ctx, user.Id, updates)
	require.Error(t, err)
	require.Empty(t, resp)

	assert.Equal(t, models.ErrInvalidArgument, errors.Unwrap(err))
}

func TestUpdateUser_NoUpdates(t *testing.T) {
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

	updates := map[string]string{
		"name":      "testuser",
	}

	mockUserProvider.EXPECT().
		UserById(ctx, int64(1)).
		Return(user, nil)

	resp, err := service.UpdateUser(ctx, user.Id, updates)
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	assert.Equal(t, user, resp)
}

func TestUpdateUser_NameIsTaken(t *testing.T) {
	service, mockUserSaver, mockUserProvider := setup.TestUser(t)

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

	updates := map[string]string{
		"name": "testuser2",
	}

	mockUserProvider.EXPECT().
		UserById(ctx, int64(1)).
		Return(user, nil)

	mockUserSaver.EXPECT().
		UpdateUser(ctx, user).
		Return(models.ErrNameIsTaken)

	resp, err := service.UpdateUser(ctx, user.Id, updates)
	require.Error(t, err)
	require.Empty(t, resp)

	assert.Equal(t, models.ErrNameIsTaken, errors.Unwrap(err))
}