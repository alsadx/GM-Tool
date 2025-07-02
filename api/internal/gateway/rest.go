package gateway

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	campaignv1 "github.com/alsadx/gm-protos/gen/go/campaignv1"
	ssov1 "github.com/alsadx/gm-protos/gen/go/ssov1"

	"gateway/internal/http/handlers"
	"gateway/internal/middleware"
)

func RunRest(campaignClient campaignv1.CampaignToolClient, userClient ssov1.UserInfoClient) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// three routers: for public routes, for private routes and for custom routes
	publicMux := runtime.NewServeMux()
	privateMux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			if claims, ok := ctx.Value(middleware.UserClaimsKey).(*middleware.Claims); ok {
				return metadata.Pairs("user_id", strconv.FormatInt(claims.UID, 10))
			}
			return nil
		}),
	)
	customMux := http.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// register public methods
	if err := ssov1.RegisterAuthHandlerFromEndpoint(ctx, publicMux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register auth handler: %v", err)
	}

	// register private methods
	if err := ssov1.RegisterUserInfoHandlerFromEndpoint(ctx, privateMux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register user info handler: %v", err)
	}
	if err := campaignv1.RegisterCampaignToolHandlerFromEndpoint(ctx, privateMux, "localhost:50052", opts); err != nil {
		log.Fatalf("failed to register campaign handler: %v", err)
	}

	// register custom methods
	customMux.HandleFunc("/api/campaigns/created", handlers.GetCreatedCampaignsHandler(campaignClient, userClient))

	// wrap privateMux with auth middleware
	protected := middleware.AuthMiddleware(privateMux)
	protectedWithCustom := middleware.AuthMiddleware(customMux)

	// main handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// TODO: придумать как по-другому разделять маршруты

		if strings.HasPrefix(path, "/api/auth/") {
			publicMux.ServeHTTP(w, r)
			return
		} else if strings.HasPrefix(path, "/api/campaigns/created") {
			protectedWithCustom.ServeHTTP(w, r)
			return
		} else {
			protected.ServeHTTP(w, r)
		}
	})

	log.Printf("REST сервер запущен на :8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		panic(err)
	}
}
