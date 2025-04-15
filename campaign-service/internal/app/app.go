package app

import (
	grpcapp "campaigntool/internal/app/grpc"
	"campaigntool/internal/config"
	"campaigntool/internal/storage/postgres"
	"campaigntool/internal/services/campaign"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
	Storage    *postgres.Storage
}

func New(log *slog.Logger, grpcPort int, dbConfig *config.DBConfig) *App {
	storage, err := postgres.New(dbConfig)
	if err != nil {
		panic(err)
	}

	campaignToolService := campaign.New(log, storage, storage)

	grpcApp := grpcapp.New(log, campaignToolService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    storage,
	}
}
