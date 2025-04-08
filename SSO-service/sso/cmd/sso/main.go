package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/hash"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting SSO server", slog.Any("config", cfg), slog.Int("grpc_port", cfg.GRPC.Port))

	hasher := hash.NewHasher()

	tokenManager := jwt.NewTokenManager()

	application := app.New(log, cfg.GRPC.Port, &cfg.DB, cfg.Auth.AccessTokenTTL, hasher, tokenManager)
	defer application.Storage.Close()

	go application.GRPCServer.MustRun()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()

	log.Info("SSO server is stopped")
}
