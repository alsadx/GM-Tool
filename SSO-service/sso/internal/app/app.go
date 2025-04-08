package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	"sso/internal/lib/hash"
	"sso/internal/lib/jwt"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
	Storage    *postgres.Storage
}

func New(log *slog.Logger, grpcPort int, dbConfig *config.DBConfig, tokenTTL time.Duration, hasher *hash.Hasher, tokenManager *jwt.TokenManager) *App {
	storage, err := postgres.New(dbConfig)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL, hasher, tokenManager)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    storage,
	}
}
