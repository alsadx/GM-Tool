package gateway

import (
	"context"
	"log"
	"net/http"
	"strconv"
	// "strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	campaignv1 "github.com/alsadx/gm-protos/gen/go/campaignv1"
	ssov1 "github.com/alsadx/gm-protos/gen/go/ssov1"

	"github.com/alsadx/GM-Tool/api/internal/http/handlers"
	"github.com/alsadx/GM-Tool/api/internal/middleware"
)

func RunRest(campaignClient campaignv1.CampaignToolClient, userClient ssov1.UserInfoClient) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
  
	// create mux for private routes
	mux := runtime.NewServeMux(
        runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
            if claims, ok := ctx.Value(middleware.UserClaimsKey).(*middleware.Claims); ok {
                return metadata.Pairs("user_id", strconv.FormatInt(claims.UID, 10))
            }
            return nil
        }),
    )

    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
  
	// register grpc services
    if err := ssov1.RegisterUserInfoHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
        log.Fatalf("failed to register user info handler: %v", err)
    }
    if err := campaignv1.RegisterCampaignToolHandlerFromEndpoint(ctx, mux, "localhost:50052", opts); err != nil {
        log.Fatalf("failed to register campaign handler: %v", err)
    }
  
	// register custom handler
    if err := mux.HandlePath(
        "GET",
        "/api/campaigns/created",
        func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
            handlers.GetCreatedCampaignsHandler(campaignClient, userClient)(w, r)
        },
    ); err != nil {
        log.Fatalf("failed to add custom handler: %v", err)
    }
  
	// wrap private routes with auth middleware
	protected := middleware.AuthMiddleware(mux)

	// mux for public routes
	publicMux := runtime.NewServeMux()
	if err := ssov1.RegisterAuthHandlerFromEndpoint(ctx, publicMux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register auth handler: %v", err)
	}
  
	// register routes
    rootMux := http.NewServeMux()
    rootMux.Handle("/api/auth/", publicMux)
    rootMux.Handle("/api/", protected)
  
  
    log.Printf("REST сервер запущен на :8081")
    if err := http.ListenAndServe(":8081", rootMux); err != nil {
        panic(err)
    }
  }

//   func RunRest(campaignClient campaignv1.CampaignToolClient, userClient ssov1.UserInfoClient) {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()
  
// 	// three routers: for public routes, for private routes and for custom routes
// 	publicMux := runtime.NewServeMux()
// 	privateMux := runtime.NewServeMux(
// 	  runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
// 		if claims, ok := ctx.Value(middleware.UserClaimsKey).(*middleware.Claims); ok {
// 		  return metadata.Pairs("user_id", strconv.FormatInt(claims.UID, 10))
// 		}
// 		return nil
// 	  }),
// 	)
  
// 	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
  
// 	// register public methods
// 	if err := ssov1.RegisterAuthHandlerFromEndpoint(ctx, publicMux, "localhost:50051", opts); err != nil {
// 	  log.Fatalf("failed to register auth handler: %v", err)
// 	}
  
// 	// register private methods
// 	if err := ssov1.RegisterUserInfoHandlerFromEndpoint(ctx, privateMux, "localhost:50051", opts); err != nil {
// 	  log.Fatalf("failed to register user info handler: %v", err)
// 	}
// 	if err := campaignv1.RegisterCampaignToolHandlerFromEndpoint(ctx, privateMux, "localhost:50052", opts); err != nil {
// 	  log.Fatalf("failed to register campaign handler: %v", err)
// 	}
  
// 	if err := privateMux.HandlePath(
// 	  "GET",
// 	  "/api/campaigns/created",
// 	  func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
// 		handlerFunc := handlers.GetCreatedCampaignsHandler(campaignClient, userClient)
// 		handlerFunc(w, r)
// 	  },
// 	); err != nil {
// 	  log.Fatalf("failed to add custom handler: %v", err)
// 	}
  
// 	// wrap privateMux with auth middleware
// 	protected := middleware.AuthMiddleware(privateMux)
  
// 	rootMux := http.NewServeMux()
	
// 	rootMux.Handle("/api/auth/", publicMux)
// 	rootMux.Handle("/api/", protected)
  
  
// 	log.Printf("REST сервер запущен на :8081")
// 	if err := http.ListenAndServe(":8081", rootMux); err != nil {
// 	  panic(err)
// 	}
//   }