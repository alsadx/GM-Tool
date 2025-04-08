package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alsadx/GM-Tool/internal/auth"
	"github.com/alsadx/GM-Tool/internal/config"
	"github.com/alsadx/GM-Tool/internal/delivery/http"
	"github.com/alsadx/GM-Tool/internal/hash"
	"github.com/alsadx/GM-Tool/internal/repository"
	"github.com/alsadx/GM-Tool/internal/repository/postgres"
	"github.com/alsadx/GM-Tool/internal/service"
	"github.com/alsadx/GM-Tool/server"
)

func Run() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Printf("%+v\n", cfg)

	dbPool, err := postgres.InitDbPool(&cfg.DB)
	if err != nil {
		log.Error("failed to init db pool: " + err.Error())
		return
	}

	hasher := hash.NewHasher()
	tokenManager, err := auth.NewManager("123")
	if err != nil {
		log.Error("failed to init token manager: " + err.Error())
		return
	}

	repositories := repository.NewRepositories(dbPool)

	serviceDeps := service.ServisesDeps{
		Repos:           repositories,
		TokenManager:    tokenManager,
		Hasher:          *hasher,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	}

	services := service.NewServices(serviceDeps)

	handler := http.NewHandler(services.Players, services.Masters, tokenManager)

	server := server.NewServer(cfg, handler.Init(cfg.HttpServer.Host, cfg.HttpServer.Port))

	if err := server.Run(); err != nil {
		log.Error("failed to run server: " + err.Error())
		return
	}

}
