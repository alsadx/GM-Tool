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
	Login(ctx context.Context, email, password string, appId int) (token string, err error)
	IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error)
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
		AppId:    req.AppId,
	}

	if err := ValidateInput(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, input.Email, input.Password, int(input.AppId))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
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
			case "appid":
				return fmt.Errorf("app_id is required")
			case "userid":
				return fmt.Errorf("user_id is required")
			}
		}
	}

	return nil
}
