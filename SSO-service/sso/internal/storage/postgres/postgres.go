package postgres

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/config"
	"sso/internal/domain/models"
	"sso/internal/storage"
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
func (s *Storage) SaveUser(ctx context.Context, email, name string, passHash string) (int64, error) {
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
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

// User returns user from storage by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	op := "storage.postgres.User"

	user := models.User{}

	query := `
		SELECT id, email, pass_hash, name
		FROM users
		WHERE email = $1
	`

	err := s.dbPool.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.PassHash, &user.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
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
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appId int) (models.App, error) {
	op := "storage.postgres.App"

	var app models.App
	query := `
		SELECT id, name, secret
		FROM apps
		WHERE id = $1
	`

	err := s.dbPool.QueryRow(ctx, query, appId).Scan(&app.Id, &app.Name, &app.SigningKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) Close() {
	s.dbPool.Close()
}
