package main

import (
	"campaigntool/internal/app"
	"campaigntool/internal/config"
	"campaigntool/internal/logger"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func logMemoryUsage(log *slog.Logger) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	log.Info("Alloc", slog.Any("MiB", memStats.Alloc/1024/1024))
	log.Info("TotalAlloc", slog.Any("MiB", memStats.TotalAlloc/1024/1024))
	log.Info("Sys", slog.Any("MiB", memStats.Sys/1024/1024))
}

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	logMemoryUsage(log)

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

	logMemoryUsage(log)
}
