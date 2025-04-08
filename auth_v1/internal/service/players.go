package service

import (
	"context"
	"time"

	"github.com/alsadx/GM-Tool/internal/auth"
	"github.com/alsadx/GM-Tool/internal/domain"
	"github.com/alsadx/GM-Tool/internal/hash"
	"github.com/alsadx/GM-Tool/internal/repository"
)

type PlayersService struct {
	repo         repository.Players
	hasher       hash.Hasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewPlayersService(repo repository.Players, hasher hash.Hasher, tokenManager auth.TokenManager, accessTTL, refreshTTL time.Duration) *PlayersService {
	return &PlayersService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *PlayersService) SignUp(ctx context.Context, input SignUpInput) error {
	player := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     s.hasher.Hash(input.Password),
		RegisteredAt: time.Now(),
	}

	if err := s.repo.CreatePlayer(ctx, player); err != nil {
		return err
	}

	return nil
}

func (s *PlayersService) SignIn(ctx context.Context, input SignInInput) (Tokens, error) {
	player, err := s.repo.GetByCred(ctx, input.Email, input.Password)
	if err != nil {
		return Tokens{}, err
	}

	return s.CreateSession(ctx, player.ID)
}

func (s *PlayersService) CreateSession(ctx context.Context, playerId string) (Tokens, error) {
	var res Tokens
	var err error

	res.AccessToken, err = s.tokenManager.NewJWT(playerId, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.accessTokenTTL),
	}

	if err := s.repo.SetSession(ctx, playerId, session); err != nil {
		return res, err
	}

	return res, nil
}

func (s *PlayersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	player, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.CreateSession(ctx, player.ID)
}