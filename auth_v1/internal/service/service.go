package service

import (
	"context"
	"time"

	"github.com/alsadx/GM-Tool/internal/auth"
	"github.com/alsadx/GM-Tool/internal/hash"
	"github.com/alsadx/GM-Tool/internal/repository"
)

type SignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=3,max=64"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=3,max=64"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Players interface {
	SignUp(ctx context.Context, input SignUpInput) error
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)
	CreateSession(ctx context.Context, playerId string) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Masters interface {
	SignUp(ctx context.Context, input SignUpInput) error
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)
	CreateSession(ctx context.Context, masterId string) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Services struct {
	Players Players
	Masters Masters
}

type ServisesDeps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	Hasher          hash.Hasher
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps ServisesDeps) *Services {
	return &Services{
		Players: NewPlayersService(deps.Repos.Players, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Masters: NewMastersService(deps.Repos.Masters, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}