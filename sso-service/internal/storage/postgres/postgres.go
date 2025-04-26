package postgres

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/config"
	"sso/internal/domain/models"
	"strings"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	dbPool *pgxpool.Pool
}

func New(dbConfig *config.DBConfig) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pool config: %s", err)
	}
	config.MaxConns = dbConfig.MaxConn

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %s", err)
	}

	return &Storage{
		dbPool: pool,
	}, nil
}

// SaveUser saves user to storage
func (s *Storage) SaveUser(ctx context.Context, email, name string, passHash []byte) (int64, error) {
	op := "storage.postgres.SaveUser"
	var userId int64

	query := `
		INSERT INTO users (email, pass_hash, name)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := s.dbPool.QueryRow(ctx, query, email, passHash, name).Scan(&userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, models.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

// UpdateUser updates user info in storage
func (s *Storage) UpdateUser(ctx context.Context, user *models.User) (error) {
	op := "storage.postgres.UpdateUser"

	query := `
		UPDATE users
		SET name = $1, full_name = $2, avatar = $3
		WHERE id = $4
	`
	_, err := s.dbPool.Exec(ctx, query, user.Name, user.FullName, user.AvatarUrl, user.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, models.ErrInvalidArgument)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

// DeleteUser deletes user from storage
func (s *Storage) DeleteUser(ctx context.Context, userId int64) (error) {
	op := "storage.postgres.DeleteUser"

	query := `
		DELETE FROM users
		WHERE id = $1
	`
	commandTag, err := s.dbPool.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if commandTag.RowsAffected() == 0 {
        return fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
    }

	return nil
}

// User returns user from storage by email
func (s *Storage) UserByEmail(ctx context.Context, email string) (models.User, error) {
	op := "storage.postgres.UserByEmail"

	user := models.User{}

	query := `
		SELECT id, email, pass_hash, name, full_name, is_admin, avatar_url
		FROM users
		WHERE email = $1
	`

	err := s.dbPool.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.PassHash, &user.Name, &user.FullName, &user.IsAdmin, &user.AvatarUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

// User returns user from storage by id
func (s *Storage) UserById(ctx context.Context, userId int64) (models.User, error) {
	op := "storage.postgres.UserById"

	user := models.User{}

	query := `
		SELECT id, email, pass_hash, name, full_name, is_admin, avatar_url
		FROM users
		WHERE id = $1
	`

	err := s.dbPool.QueryRow(ctx, query, userId).Scan(&user.Id, &user.Email, &user.PassHash, &user.Name, &user.FullName, &user.IsAdmin, &user.AvatarUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) UserByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	op := "storage.postgres.UserByRefreshToken"

	var user models.User

	query := `
		SELECT u.id, u.email, u.pass_hash, u.name, u.full_name, u.is_admin, u.avatar_url
		FROM users u
		JOIN sessions s ON u.id = s.user_id
		WHERE s.refresh_token = $1 AND s.expires_at > NOW()
	`

	err := s.dbPool.QueryRow(ctx, query, refreshToken).Scan(&user.Id, &user.Email, &user.PassHash, &user.Name, &user.FullName, &user.IsAdmin, &user.AvatarUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

// IsAdmin returns true if user is admin and false otherwise
func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	op := "storage.postgres.IsAdmin"

	var isAdmin bool

	query := `
		SELECT is_admin
		FROM users
		WHERE id = $1
	`

	err := s.dbPool.QueryRow(ctx, query, userId).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) SetSession(ctx context.Context, userId int64, session models.Session) error {
	op := "storage.postgres.SetSession"

	query := `
		INSERT INTO sessions (user_id, refresh_token, expires_at)
		VALUES ($3, $1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET refresh_token = EXCLUDED.refresh_token,
		expires_at = EXCLUDED.expires_at;
	`

	_, err := s.dbPool.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSession(ctx context.Context, userId int64) error {
	op := "storage.postgres.DeleteSession"

	query := `
		DELETE FROM sessions
		WHERE user_id = $1;
	`

	commandTag, err := s.dbPool.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if commandTag.RowsAffected() == 0 {
        return fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
    }

	return nil
}

func (s *Storage) HealthCheck(ctx context.Context) error {
	return s.dbPool.Ping(ctx)
}

func (s *Storage) Close() {
	s.dbPool.Close()
}
