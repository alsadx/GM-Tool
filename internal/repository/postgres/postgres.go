package postgres

import (
	"context"
	"fmt"

	"github.com/alsadx/GM-Tool/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepository struct {
	DbPool *pgxpool.Pool
}

func InitDbPool(info *config.DbInfo) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", info.User, info.Password, info.Host, info.Port, info.Name)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pool config: %s", err)
	}
	config.MaxConns = info.MaxConn

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %s", err)
	}

	return pool, nil
}
