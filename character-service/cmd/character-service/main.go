package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	character_core "github.com/alsadx/GM-Tool/character-service/gen/service/character-core"
	character_health "github.com/alsadx/GM-Tool/character-service/gen/service/character-health"
	character_hitdice "github.com/alsadx/GM-Tool/character-service/gen/service/character-hitdice"
	character_progress "github.com/alsadx/GM-Tool/character-service/gen/service/character-progress"
	character_stats "github.com/alsadx/GM-Tool/character-service/gen/service/character-stats"
	core_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-core"
	health_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-health"
	hitdice_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-hitdice"
	progress_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-progress"
	stats_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-stats"
	handler_core "github.com/alsadx/GM-Tool/character-service/internal/handlers/character-core"
	handler_health "github.com/alsadx/GM-Tool/character-service/internal/handlers/character-health"
	handler_hitdice "github.com/alsadx/GM-Tool/character-service/internal/handlers/character-hitdice"
	handler_progress "github.com/alsadx/GM-Tool/character-service/internal/handlers/character-progress"
	handler_stats "github.com/alsadx/GM-Tool/character-service/internal/handlers/character-stats"
	mongoRepo "github.com/alsadx/GM-Tool/character-service/internal/repository/mongo"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	port := flag.Int("p", 8080, "API handler port. Can be overridden by PORT env var.")
	flag.Parse()

	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil {
			*port = p
		}
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	mongoURL := getEnv("MONGO_URL", "mongodb://localhost:27017")
	mongoDBName := getEnv("MONGO_DATABASE", "character-service-MongoDB")

	sink := zapr.NewLogger(logger).GetSink()
	loggerOptions := options.
		Logger().
		SetSink(sink).
		SetComponentLevel(options.LogComponentAll, options.LogLevelInfo)

	client, err := mongo.Connect(
		options.Client().
			ApplyURI(mongoURL).
			SetLoggerOptions(loggerOptions),
	)
	if err != nil {
		logger.Fatal("error when connecting to MongoDB", zap.Error(err))
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		logger.Fatal("MongoDB ping failed", zap.Error(err))
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	repo := mongoRepo.NewRepository(client.Database(mongoDBName))

	logger.Info("service starting...",
		zap.Int("port", *port),
		zap.String("mongo_url", mongoURL),
		zap.String("mongo_database", mongoDBName),
	)

	ctrl_core := core_controller.New(repo)
	ctrl_health := health_controller.New(repo)
	ctrl_progress := progress_controller.New(repo)
	ctrl_stats := stats_controller.New(repo)
	ctrl_hitdice := hitdice_controller.New(repo)

	handler_core := handler_core.New(ctrl_core)
	handler_health := handler_health.New(ctrl_health)
	handler_progress := handler_progress.New(ctrl_progress)
	handler_stats := handler_stats.New(ctrl_stats)
	handler_hitdice := handler_hitdice.New(ctrl_hitdice)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()

	character_core.RegisterCharacterCoreServiceServer(srv, handler_core)
	character_health.RegisterCharacterHealthServiceServer(srv, handler_health)
	character_progress.RegisterCharacterProgressServiceServer(srv, handler_progress)
	character_stats.RegisterCharacterStatsServiceServer(srv, handler_stats)
	character_hitdice.RegisterCharacterHitDiceServiceServer(srv, handler_hitdice)

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}