package auth

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/services/auth"
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
	Logout(ctx context.Context, token string) (err error)
	GetCurrentUser(ctx context.Context, token string) (user models.User, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func RegisterServerAPI(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	input := models.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := ValidateInput(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokens, err := s.auth.Login(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{Token: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	input := models.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := ValidateInput(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := s.auth.Register(ctx, input.Email, input.Password, input.Name)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{UserId: userId}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
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

	isAdmin, err := s.auth.IsAdmin(ctx, input.UserId)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func (s *serverAPI) RefreshToken(ctx context.Context, req *ssov1.RefreshTokenRequest) (*ssov1.RefreshTokenResponse, error) {
	tokens, err := s.auth.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid or expired refresh token")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.RefreshTokenResponse{Token: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	err := s.auth.Logout(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return &ssov1.LogoutResponse{Success: false}, status.Errorf(codes.NotFound, "user not found")
		}
		return &ssov1.LogoutResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.LogoutResponse{Success: true}, nil
}

func (s *serverAPI) GetCurrentUser(ctx context.Context, req *ssov1.GetCurrentUserRequest) (*ssov1.GetCurrentUserResponse, error) {
	user, err := s.auth.GetCurrentUser(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.GetCurrentUserResponse{UserId: user.Id, Name: user.Name, Email: user.Email}, nil
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
