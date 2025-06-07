package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	character_core "github.com/alsadx/GM-Tool/character-service/gen/service/character-core"
	character_health "github.com/alsadx/GM-Tool/character-service/gen/service/character-health"
	core_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-core"
	health_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-health"
	handler_core "github.com/alsadx/GM-Tool/character-service/internal/handlers/character-core"
	handler_health "github.com/alsadx/GM-Tool/character-service/internal/handlers/character-health"
	mongoRepo "github.com/alsadx/GM-Tool/character-service/internal/repository/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	port := flag.Int("p", 0, "API handler port")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	
	defer logger.Sync()

	urlMongo := "mongodb://localhost:27017"
	client, err := mongo.Connect(options.Client().ApplyURI(urlMongo))
	if err != nil {
		logger.Fatal("error when connecting to MongoDB", zap.Error(err))
		panic(err)
	}
	defer func () {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
	repo := mongoRepo.NewRepository(client.Database("character-service-MongoDB"))

	ctrl_core := core_controller.New(repo)
	ctrl_health := health_controller.New(repo)

	handler_core := handler_core.New(ctrl_core)
	handler_health := handler_health.New(ctrl_health)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	
	character_core.RegisterCharacterCoreServiceServer(srv, handler_core)
	character_health.RegisterCharacterHealthServiceServer(srv, handler_health)

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}