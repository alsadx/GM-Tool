package main

import (
	"fmt"
	"gateway/internal/gateway"
	"log"

	"github.com/alsadx/gm-protos/gen/go/campaignv1"
	"github.com/alsadx/gm-protos/gen/go/ssov1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create SSO clients
	cc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc server connection failed: %v", err)
	}
	// authClient := ssov1.NewAuthClient(cc)
	// fmt.Println("Auth client created")

	userInfoClient := ssov1.NewUserInfoClient(cc)
	fmt.Println("Auth client created")

	// Create campaign client
	cc, err = grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc server connection failed: %v", err)
	}
	campaignClient := campaignv1.NewCampaignToolClient(cc)
	fmt.Println("Campaign client created")

	gateway.RunRest(campaignClient, userInfoClient)
}
