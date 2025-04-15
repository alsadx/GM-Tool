package main

import (
	"campaigntool/internal/app"
	"campaigntool/internal/config"
	"campaigntool/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting CampaignTool server", slog.Any("config", cfg), slog.Int("grpc_port", cfg.GRPC.Port))

	application := app.New(log, cfg.GRPC.Port, &cfg.DB)
	defer application.Storage.Close()

	go application.GRPCServer.MustRun()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()

	log.Info("CampaignTool server is stopped")
}
