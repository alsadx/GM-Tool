package auth

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"strings"

	ssov1 "github.com/alsadx/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, email, password, name string) (userId int64, err error)
	Login(ctx context.Context, email, password string) (token models.Tokens, err error)
	IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error)
	RefreshToken(ctx context.Context, refreshToken string) (token models.Tokens, err error)
	Logout(ctx context.Context, userId int64) (err error)
	GetCurrentUser(ctx context.Context, token string) (user models.User, err error)
	HealthCheck(ctx context.Context) (err error)
}

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	Auth Auth
}

func RegisterServerAPI(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &ServerAPI{Auth: auth})
}

// func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
// 	}

// 	authHeader, ok := md["authorization"]
// 	if !ok || len(authHeader) == 0 {
// 		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
// 	}

// 	token, err := jwt.ExtractTokenFromHeader(authHeader[0])
// 	if err != nil {
// 		return nil, err
// 	}

// 	// if os.Getenv("TEST_ENV") == "true" {
// 	// 	if token == "valid-token" {
// 	// 		ctx = context.WithValue(ctx, "user_id", 1)
// 	// 		return handler(ctx, req)
// 	// 	} else {
// 	// 		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
// 	// 	}
// 	// }

// 	userId, err := jwt.ValidateToken(token, "secret")
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
// 	}

// 	// if err := jwt.IsTokenExpired(token); err != nil {
//     //     return nil, status.Errorf(codes.Unauthenticated, "token expired")
//     // }

// 	ctx = context.WithValue(ctx, "user_id", userId)

// 	return handler(ctx, req)
// }

func (s *ServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	input := models.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := ValidateInput(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokens, err := s.Auth.Login(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{Token: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	input := models.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := ValidateInput(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := s.Auth.Register(ctx, input.Email, input.Password, input.Name)
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{UserId: userId}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	input := models.IsAdminInput{
		UserId: req.UserId,
	}

	// validate := validator.New()
	// if err := validate.Struct(input); err != nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	// }

	if err := ValidateInput(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	isAdmin, err := s.Auth.IsAdmin(ctx, input.UserId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func (s *ServerAPI) RefreshToken(ctx context.Context, req *ssov1.RefreshTokenRequest) (*ssov1.RefreshTokenResponse, error) {
	tokens, err := s.Auth.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, models.ErrInvalidRefreshToken) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid or expired refresh token")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.RefreshTokenResponse{Token: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *ServerAPI) Logout(ctx context.Context, req *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return &ssov1.LogoutResponse{Success: false}, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.Auth.Logout(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return &ssov1.LogoutResponse{Success: false}, status.Errorf(codes.NotFound, "user not found")
		}
		return &ssov1.LogoutResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.LogoutResponse{Success: true}, nil
}

func (s *ServerAPI) GetCurrentUser(ctx context.Context, req *ssov1.GetCurrentUserRequest) (*ssov1.GetCurrentUserResponse, error) {
	user, err := s.Auth.GetCurrentUser(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.GetCurrentUserResponse{UserId: user.Id, Name: user.Name, Email: user.Email}, nil
}

func (s *ServerAPI) HealthCheck(ctx context.Context, req *ssov1.HealthCheckRequest) (*ssov1.HealthCheckResponse, error) {
	// Проверяем подключение к базе данных
	err := s.Auth.HealthCheck(ctx)
	if err != nil {
		return &ssov1.HealthCheckResponse{Status: ssov1.HealthCheckResponse_NOT_SERVING}, nil
	}
	return &ssov1.HealthCheckResponse{Status: ssov1.HealthCheckResponse_SERVING}, nil
}

func ValidateInput(input any) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(input)
	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			fmt.Printf("field: %s\n", field)
			switch field {
			case "email":
				return fmt.Errorf("email is required")
			case "name":
				return fmt.Errorf("name is required")
			case "password":
				return fmt.Errorf("password is required")
			case "userid":
				return fmt.Errorf("user_id is required")
			}
		}
	}

	return nil
}
