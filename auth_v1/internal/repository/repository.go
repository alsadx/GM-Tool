package repository

import (
	"context"

	"github.com/alsadx/GM-Tool/internal/domain"
	"github.com/alsadx/GM-Tool/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Players interface {
	CreatePlayer(ctx context.Context, user domain.User) error
	GetByCred(ctx context.Context, email string, password string) (domain.User, error)
    GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
    GetByID(ctx context.Context, id string) (domain.User, error)
    SetSession(ctx context.Context, playerID string, session domain.Session) error
}
type Masters interface {
	CreateMaster(ctx context.Context, user domain.User) error
    GetByCred(ctx context.Context, email string, password string) (domain.User, error)
    GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
    GetByID(ctx context.Context, id string) (domain.User, error)
    SetSession(ctx context.Context, masterID string, session domain.Session) error
}

type Repositories struct {
	Players Players
	Masters Masters
}

func NewRepositories(dbPool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Players: postgres.NewPlayersRepo(dbPool),
		Masters: postgres.NewMastersRepo(dbPool),
	}
}
