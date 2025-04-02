package postgres

import (
	"context"
	"fmt"

	"github.com/alsadx/GM-Tool/internal/domain"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MastersRepo struct {
	BaseRepository
}

func NewMastersRepo(dbPool *pgxpool.Pool) *MastersRepo {
	return &MastersRepo{
		BaseRepository: BaseRepository{DbPool: dbPool},
	}
}

func (r *MastersRepo) CreateMaster(ctx context.Context, user domain.User) error {
	query := `
        INSERT INTO masters (username, email, password_hash, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	err := r.DbPool.QueryRow(ctx, query, user.Name, user.Email, user.Password, user.RegisteredAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create master: %s", err)
	}

	return nil
}

func (r *MastersRepo) GetByCred(ctx context.Context, email string, password string) (domain.User, error) {
	query := `
        SELECT id, userame, email, password, created_at, session_refresh_token, session_expires_at
        FROM masters
        WHERE email = $1 AND password = $2
    `
	var master domain.User
	err := r.DbPool.QueryRow(ctx, query, email, password).Scan(
		&master.ID,
		&master.Name,
		&master.Email,
		&master.Password,
		&master.RegisteredAt,
		&master.Session.RefreshToken,
		&master.Session.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("master not found")
		}
		return domain.User{}, fmt.Errorf("failed to get master by credentials: %s", err)
	}

	return master, nil
}

func (r *MastersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	query := `
		SELECT id, userame, email, password, created_at, session_refresh_token, session_expires_at
		FROM masters
		WHERE session_refresh_token = $1 AND session_expires_at > NOW()
	`

	var master domain.User
	err := r.DbPool.QueryRow(ctx, query, refreshToken).Scan(
		&master.ID,
		&master.Name,
		&master.Email,
		&master.Password,
		&master.RegisteredAt,
		&master.Session.RefreshToken,
		&master.Session.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("master not found")
		}
		return domain.User{}, fmt.Errorf("failed to get master by refresh token: %s", err)
	}

	return master, nil
}

func (r *MastersRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	query := `
		SELECT id, userame, email, password, created_at, session_refresh_token, session_expires_at
		FROM masters
		WHERE id = $1 AND verified = true
	`
	var master domain.User
	err := r.DbPool.QueryRow(ctx, query, id).Scan(
		&master.ID,
		&master.Name,
		&master.Email,
		&master.Password,
		&master.RegisteredAt,
		&master.Session.RefreshToken,
		&master.Session.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("master not found")
		}
		return domain.User{}, fmt.Errorf("failed to get master by id: %s", err)
	}

	return master, nil
}

func (r *MastersRepo) SetSession(ctx context.Context, masterID string, session domain.Session) error {
	query := `
	UPDATE masters
	SET session_refresh_token = $1, session_expires_at = $2
	WHERE id = $3
`

	_, err := r.DbPool.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, masterID)
	if err != nil {
		return fmt.Errorf("failed to set session: %s", err)
	}

	return nil
}
