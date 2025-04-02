package postgres

import (
	"context"
	"fmt"

	"github.com/alsadx/GM-Tool/internal/domain"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayersRepo struct {
	BaseRepository
}

func NewPlayersRepo(dbPool *pgxpool.Pool) *PlayersRepo {
	return &PlayersRepo{
		BaseRepository: BaseRepository{DbPool: dbPool},
	}
}

func (r *PlayersRepo) CreatePlayer(ctx context.Context, user domain.User) error {
	query := `
        INSERT INTO players (username, email, password_hash, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	err := r.DbPool.QueryRow(ctx, query, user.Name, user.Email, user.Password, user.RegisteredAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create player: %s", err)
	}

	return nil
}

func (r *PlayersRepo) GetByCred(ctx context.Context, email string, password string) (domain.User, error) {
	query := `
        SELECT id, userame, email, password, created_at, session_refresh_token, session_expires_at
        FROM players
        WHERE email = $1 AND password = $2
    `
	var player domain.User
	err := r.DbPool.QueryRow(ctx, query, email, password).Scan(
		&player.ID,
		&player.Name,
		&player.Email,
		&player.Password,
		&player.RegisteredAt,
		&player.Session.RefreshToken,
		&player.Session.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("player not found")
		}
		return domain.User{}, fmt.Errorf("failed to get player by credentials: %s", err)
	}

	return player, nil
}

func (r *PlayersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	query := `
		SELECT id, userame, email, password, created_at, session_refresh_token, session_expires_at
		FROM players
		WHERE session_refresh_token = $1 AND session_expires_at > NOW()
	`

	var player domain.User
	err := r.DbPool.QueryRow(ctx, query, refreshToken).Scan(
		&player.ID,
		&player.Name,
		&player.Email,
		&player.Password,
		&player.RegisteredAt,
		&player.Session.RefreshToken,
		&player.Session.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("player not found")
		}
		return domain.User{}, fmt.Errorf("failed to get player by refresh token: %s", err)
	}

	return player, nil
}

func (r *PlayersRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	query := `
		SELECT id, userame, email, password, created_at, session_refresh_token, session_expires_at
		FROM players
		WHERE id = $1 AND verified = true
	`
	var player domain.User
	err := r.DbPool.QueryRow(ctx, query, id).Scan(
		&player.ID,
		&player.Name,
		&player.Email,
		&player.Password,
		&player.RegisteredAt,
		&player.Session.RefreshToken,
		&player.Session.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("player not found")
		}
		return domain.User{}, fmt.Errorf("failed to get player by id: %s", err)
	}

	return player, nil
}

func (r *PlayersRepo) SetSession(ctx context.Context, playerID string, session domain.Session) error {
	query := `
	UPDATE players
	SET session_refresh_token = $1, session_expires_at = $2
	WHERE id = $3
`

	_, err := r.DbPool.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, playerID)
	if err != nil {
		return fmt.Errorf("failed to set session: %s", err)
	}

	return nil
}
