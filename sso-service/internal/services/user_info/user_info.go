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
	log          *slog.Logger
	userSaver    srv.UserSaver
	userProvider srv.UserProvider
}

func New(log *slog.Logger, userSaver srv.UserSaver, userProvider srv.UserProvider) *UserInfo {
	return &UserInfo{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
	}
}

func (u *UserInfo) GetUserById(ctx context.Context, userId int64) (*models.User, error) {
	const op = "auth.GetUserById"

	log := u.log.With(slog.String("op", op))

	log.Info("getting user by id")

	user, err := u.userProvider.UserById(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.log.Warn("user not found", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.log.Error("failed to get user by id", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got user by id", slog.Int64("id", user.Id))

	return user, nil
}

func (u *UserInfo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "auth.GetUserByEmail"

	log := u.log.With(slog.String("op", op))

	log.Info("getting user by email")

	user, err := u.userProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.log.Warn("user not found", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.log.Error("failed to get user by email", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got user by email", slog.String("email", user.Email))

	return user, nil
}
func (u *UserInfo) UpdateUser(ctx context.Context, userId int64, updates map[string]string) (*models.User, error) {
	const op = "storage.postgres.UpdateUser"

	user, err := u.userProvider.UserById(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.log.Warn("user not found", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.log.Error("failed to get user by id", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	for key, value := range updates {
		switch key {
		case "name":
			user.Name = value
		case "full_name":
			user.FullName = value
		case "avatar":
			user.AvatarUrl = value
		}
	}

	err = u.userSaver.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrInvalidArgument) {
			u.log.Warn("name is taken", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrInvalidArgument)
		}
		u.log.Error("failed to update user", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *UserInfo) DeleteUser(ctx context.Context, userId int64) error {
	const op = "storage.postgres.DeleteUser"

	err := u.userSaver.DeleteSession(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.log.Warn("user not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.log.Error("failed to delete session", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	// TODO: delete user from campaigns and players

	err = u.userSaver.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			u.log.Warn("user not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		u.log.Error("failed to delete user", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
