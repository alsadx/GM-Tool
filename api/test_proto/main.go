package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alsadx/gm-protos/gen/go/ssov1"
	"github.com/alsadx/gm-protos/gen/go/campaignv1"
)

func runRest() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := ssov1.RegisterAuthHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
        log.Fatalf("failed to register auth handler: %v", err)
    }

	if err := ssov1.RegisterUserInfoHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
        log.Fatalf("failed to register user info handler: %v", err)
    }

    if err := campaignv1.RegisterCampaignToolHandlerFromEndpoint(ctx, mux, "localhost:50052", opts); err != nil {
        log.Fatalf("failed to register campaign handler: %v", err)
    }

	log.Printf("server listening at 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		panic(err)
	}
}

func main() {
	runRest()
}
