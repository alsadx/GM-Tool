package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"sso/internal/domain/models"
	"sso/internal/lib/hash"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email, name string, passHash string) (userId int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
	hasher       *hash.Hasher
	tokenManager *jwt.TokenManager
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAppNotFound        = errors.New("app not found")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration, hasher *hash.Hasher, tokenManager *jwt.TokenManager) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

// Login checks if user with given credentials exists
// if user exists returns token
// if user doesn't exist returns error
func (a *Auth) Login(ctx context.Context, email, password string, appId int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("attempting to login")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := a.hasher.CheckPassword(password, user.PassHash); err != nil {
		a.log.Warn("invalid password", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}
		a.log.Error("error getting app", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := a.tokenManager.NewJWT(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in")

	return token, nil
}

// RegisterNewUser registers new user and returns user id
// returns error if user already exists
func (a *Auth) Register(ctx context.Context, email, password, name string) (int64, error) {
	const op = "auth.Register"
	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("registering new user")

	passHash, err := a.hasher.Hash(password)
	if err != nil {
		a.log.Error("failed to generate password hash", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userId, err := a.userSaver.SaveUser(ctx, email, name, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("user already exists", slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		a.log.Error("failed to save user", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered", slog.Int64("userId", userId))

	return userId, nil
}

// IsAdmin returns true if user is admin
func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(slog.String("op", op), slog.Int64("userId", userId))

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", slog.String("error", err.Error()))

			return false, fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}
		a.log.Error("failed to check if user is admin", slog.String("error", err.Error()))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
