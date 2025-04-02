package service

import (
	"context"
	"time"

	"github.com/alsadx/GM-Tool/internal/auth"
	"github.com/alsadx/GM-Tool/internal/domain"
	"github.com/alsadx/GM-Tool/internal/hash"
	"github.com/alsadx/GM-Tool/internal/repository"
)

type MastersService struct {
	repo         repository.Masters
	hasher       hash.Hasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewMastersService(repo repository.Masters, hasher hash.Hasher, tokenManager auth.TokenManager, accessTTL, refreshTTL time.Duration) *MastersService {
	return &MastersService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *MastersService) SignUp(ctx context.Context, input SignUpInput) error {
	master := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     s.hasher.Hash(input.Password),
		RegisteredAt: time.Now(),
	}

	if err := s.repo.CreateMaster(ctx, master); err != nil {
		return err
	}

	return nil
}

func (s *MastersService) SignIn(ctx context.Context, input SignInInput) (Tokens, error) {
	master, err := s.repo.GetByCred(ctx, input.Email, input.Password)
	if err != nil {
		return Tokens{}, err
	}

	return s.CreateSession(ctx, master.ID)
}

func (s *MastersService) CreateSession(ctx context.Context, masterId string) (Tokens, error) {
	var res Tokens
	var err error

	res.AccessToken, err = s.tokenManager.NewJWT(masterId, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	if err := s.repo.SetSession(ctx, masterId, session); err != nil {
		return res, err
	}

	return res, nil
}

func (s *MastersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	master, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.CreateSession(ctx, master.ID)
}

// Получаем хэш и роль
// err := pool.QueryRow(context.Background(), `
//     SELECT id, role, password_hash FROM auth_users WHERE username = $1
// `, username).Scan(&userID, &role, &storedHash)
// if err != nil {
//     return 0, "", fmt.Errorf("failed to fetch user: %w", err)
// }

// Проверяем пароль
// err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
// if err != nil {
//     if err == bcrypt.ErrMismatchedHashAndPassword {
//         return 0, "", fmt.Errorf("invalid password")
//     }
//     return 0, "", fmt.Errorf("failed to compare passwords: %w", err)
// }
