package userinfo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"sso/internal/domain/models"
	srv "sso/internal/services"
)

type UserInfo struct {
	Log          *slog.Logger
	UserSaver    srv.UserSaver
	UserProvider srv.UserProvider
}

func New(log *slog.Logger, userSaver srv.UserSaver, userProvider srv.UserProvider) *UserInfo {
	return &UserInfo{
		Log:          log,
		UserSaver:    userSaver,
		UserProvider: userProvider,
	}
}

func (u *UserInfo) GetUserById(ctx context.Context, userId int64) (*models.User, error) {
	const op = "user_info.GetUserById"

	log := u.Log.With(slog.String("op", op))

	log.Info("getting user by id")

	user, err := u.UserProvider.UserById(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.Log.Warn("user not found", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.Log.Error("failed to get user by id", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got user by id", slog.Int64("id", user.Id))

	return user, nil
}

func (u *UserInfo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "user_info.GetUserByEmail"

	log := u.Log.With(slog.String("op", op))

	log.Info("getting user by email")

	user, err := u.UserProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.Log.Warn("user not found", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.Log.Error("failed to get user by email", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got user by email", slog.String("email", user.Email))

	return user, nil
}
func (u *UserInfo) UpdateUser(ctx context.Context, userId int64, updates map[string]string) (*models.User, error) {
	const op = "user_info.UpdateUser"

	user, err := u.UserProvider.UserById(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.Log.Warn("user not found", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.Log.Error("failed to get user by id", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	if len(updates) == 0 {
        u.Log.Info("no updates provided for the user", slog.Int64("user_id", userId))

        return nil, fmt.Errorf("%s: %w", op, models.ErrInvalidArgument)
    }

	hasChanges := false

	// TODO: validate updates

	for key, value := range updates {
		if value == "" || len(value) > 255 {
			u.Log.Warn("invalid update value", slog.String("key", key), slog.String("value", value))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrInvalidArgument)
		}
		switch key {
		case "name":
			if user.Name != value {
				user.Name = value
				hasChanges = true
			}
		case "full_name":
			if user.FullName != value {
				user.FullName = value
				hasChanges = true
			}
		case "avatar_url":
			if user.AvatarUrl != value {
				user.AvatarUrl = value
				hasChanges = true
			}
		default:
			u.Log.Warn("unknown update key", slog.String("key", key))
		}
	}

	if !hasChanges {
		u.Log.Info("no changes detected for the user", slog.Int64("user_id", userId))

		return user, nil
	}

	err = u.UserSaver.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrNameIsTaken) {
			u.Log.Warn("name is taken", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrNameIsTaken)
		}
		u.Log.Error("failed to update user", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *UserInfo) DeleteUser(ctx context.Context, userId int64) error {
	const op = "user_info.DeleteUser"

	err := u.UserSaver.DeleteSession(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.Log.Warn("user not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.Log.Error("failed to delete session", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.UserSaver.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.Log.Warn("user not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.Log.Error("failed to delete user", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
