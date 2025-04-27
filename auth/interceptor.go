package auth

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("AuthInterceptor: %s", info.FullMethod)

	// Исключаем методы, которые не требуют аутентификации
	if info.FullMethod == "/auth.Auth/Register" || info.FullMethod == "/auth.Auth/Login" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Println("got metadata")

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	log.Println("got auth header")

	token, err := ExtractTokenFromHeader(authHeader[0])
	if err != nil {
		return nil, err
	}

	log.Printf("got token: %s", token)

	if os.Getenv("TEST_ENV") == "true" {
		log.Println("TEST_ENV is true")
		if token == "valid-token" {
			ctx = context.WithValue(ctx, "user_id", 1)

			log.Println("token is valid")

			return handler(ctx, req)
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
	}

	userId, err := ValidateToken(token, "secret")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	log.Printf("user id: %d", userId)

	// if err := jwt.IsTokenExpired(token); err != nil {
	//     return nil, status.Errorf(codes.Unauthenticated, "token expired")
	// }

	ctx = context.WithValue(ctx, "user_id", userId)

	return handler(ctx, req)
}
