package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/alsadx/GM-Tool/character-service/gen/service"
	controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character"
	grpchandler "github.com/alsadx/GM-Tool/character-service/internal/handler/grpc"
	mongoRepo "github.com/alsadx/GM-Tool/character-service/internal/repository/mongo"
	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	port := flag.Int("p", 0, "API handler port")
	flag.Parse()

	urlMongo := "mongodb://localhost:27017"
	client, err := mongo.Connect(options.Client().ApplyURI(urlMongo))
	if err != nil {
		panic(err)
	}
	defer func () {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
	repo := mongoRepo.NewRepository(client.Database("character-service-MongoDB"))
	ctrl := controller.New(repo)
	h := grpchandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	service.RegisterCharacterServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}