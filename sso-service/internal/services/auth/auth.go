package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	srv "sso/internal/services"
)

const (
	signKey = "secret"
)

type Hasher interface {
	Hash(password string) ([]byte, error)
	CheckPassword(password string, hash []byte) error
}

type TokenManager interface {
	NewJWT(user *models.User, signKey string, ttl time.Duration) (string, error)
	NewRefreshToken() (string, error)
	ParseJWT(token string, signKey string) (*jwt.ParsedJWT, error)
}

type Auth struct {
	Log          *slog.Logger
	UserSaver    srv.UserSaver
	UserProvider srv.UserProvider
	TokenTTL     time.Duration
	Hasher       Hasher
	TokenManager TokenManager
}

func New(log *slog.Logger, userSaver srv.UserSaver, userProvider srv.UserProvider, tokenTTL time.Duration, hasher Hasher, tokenManager TokenManager) *Auth {
	return &Auth{
		Log:          log,
		UserSaver:    userSaver,
		UserProvider: userProvider,
		TokenTTL:     tokenTTL,
		Hasher:       hasher,
		TokenManager: tokenManager,
	}
}

// Login checks if user with given credentials exists
// if user exists returns token
// if user doesn't exist returns error
func (a *Auth) Login(ctx context.Context, email, password string) (models.Tokens, error) {
	const op = "auth.Login"

	log := a.Log.With(slog.String("op", op), slog.String("email", email))

	log.Info("attempting to login")

	user, err := a.UserProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			a.Log.Warn("user not found", slog.String("error", err.Error()))

			return models.Tokens{}, fmt.Errorf("%s: %w", op, models.ErrInvalidCredentials)
		}

		a.Log.Error("failed to get user", slog.String("error", err.Error()))

		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := a.Hasher.CheckPassword(password, user.PassHash); err != nil {
		a.Log.Warn("invalid password", slog.String("error", err.Error()))

		return models.Tokens{}, fmt.Errorf("%s: %w", op, models.ErrInvalidCredentials)
	}

	tokens, err := a.CreateSession(ctx, user)
	if err != nil {
		a.Log.Error("failed to create session", slog.String("error", err.Error()))

		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in")

	return tokens, nil
}

// RegisterNewUser registers new user and returns user id
// returns error if user already exists
func (a *Auth) Register(ctx context.Context, email, password, name string) (int64, error) {
	const op = "auth.Register"

	log := a.Log.With(slog.String("op", op), slog.String("email", email))

	log.Info("registering new user")

	passHash, err := a.Hasher.Hash(password)
	if err != nil {
		a.Log.Error("failed to generate password hash", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userId, err := a.UserSaver.SaveUser(ctx, email, name, passHash)
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			a.Log.Warn("user already exists", slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s: %w", op, models.ErrUserExists)
		}
		a.Log.Error("failed to save user", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered", slog.Int64("userId", userId))

	return userId, nil
}

// IsAdmin returns true if user is admin
func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.Log.With(slog.String("op", op), slog.Int64("userId", userId))

	log.Info("checking if user is admin")

	isAdmin, err := a.UserProvider.IsAdmin(ctx, userId)
	if err != nil {
		a.Log.Error("failed to check if user is admin", slog.String("error", err.Error()))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}

func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (models.Tokens, error) {
	const op = "auth.RefreshToken"

	log := a.Log.With(slog.String("op", op))

	log.Info("refreshing token")

	user, err := a.UserProvider.UserByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			a.Log.Warn("invalid refresh token", slog.String("error", err.Error()))

			return models.Tokens{}, fmt.Errorf("%s: %w", op, models.ErrInvalidRefreshToken)
		}
		a.Log.Error("failed to get user by refresh token", slog.String("error", err.Error()))

		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := a.CreateSession(ctx, user)
	if err != nil {
		a.Log.Error("failed to create session", slog.String("error", err.Error()))

		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user refreshed token")

	return tokens, nil
}

func (a *Auth) Logout(ctx context.Context, userId int64) error {
	const op = "auth.Logout"

	log := a.Log.With(slog.String("op", op))

	log.Info("logging out")

	err := a.UserSaver.DeleteSession(ctx, userId)
	if err != nil {
		a.Log.Error("failed to delete session", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged out")

	return nil
}

func (a *Auth) GetCurrentUser(ctx context.Context, token string) (*models.User, error) {
	const op = "auth.GetCurrentUser"

	log := a.Log.With(slog.String("op", op))

	log.Info("getting current user")

	parsedToken, err := a.TokenManager.ParseJWT(token, signKey)
	if err != nil {
		a.Log.Warn("invalid token", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.UserProvider.UserByEmail(ctx, parsedToken.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			a.Log.Warn("user not found", slog.String("error", err.Error()))

			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		a.Log.Error("failed to get user", slog.String("error", err.Error()))

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got current user", slog.String("email", user.Email))

	return user, nil
}

func (a *Auth) CreateSession(ctx context.Context, user *models.User) (models.Tokens, error) {
	const op = "auth.CreateSession"

	log := a.Log.With(slog.String("op", op), slog.String("email", user.Email))

	log.Info("creating session")

	var res models.Tokens
	var err error

	res.AccessToken, err = a.TokenManager.NewJWT(user, signKey, a.TokenTTL)
	if err != nil {
		a.Log.Error("failed to generate access token", slog.String("error", err.Error()))

		return res, fmt.Errorf("%s: %w", op, err)
	}

	res.RefreshToken, err = a.TokenManager.NewRefreshToken()
	if err != nil {
		a.Log.Error("failed to generate refresh token", slog.String("error", err.Error()))

		return res, fmt.Errorf("%s: %w", op, err)
	}

	session := models.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(a.TokenTTL),
	}

	if err := a.UserSaver.SetSession(ctx, user.Id, session); err != nil {
		a.Log.Error("failed to set session", slog.String("error", err.Error()))

		return res, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("session created")

	return res, nil
}

func (a *Auth) HealthCheck(ctx context.Context) error {
	const op = "auth.HealthCheck"

	log := a.Log.With(slog.String("op", op))

	log.Info("health checking")

	err := a.UserProvider.HealthCheck(ctx)
	if err != nil {
		a.Log.Error("failed to health check", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("health checked")

	return nil
}
