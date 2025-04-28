package userinfo

import (
	"context"
	"errors"
	"sso/internal/domain/models"

	"protos/gen/go/ssov1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserInfo interface {
	GetUserById(ctx context.Context, userId int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, userId int64, updates map[string]string) (*models.User, error)
	DeleteUser(ctx context.Context, userId int64) error
}

type UserInfoAPI struct {
	ssov1.UnimplementedUserInfoServer
	UserInfo UserInfo
}

func RegisterUserInfoAPI(gRPC *grpc.Server, userInfo UserInfo) {
	ssov1.RegisterUserInfoServer(gRPC, &UserInfoAPI{UserInfo: userInfo})
}

func (s *UserInfoAPI) GetUserById(ctx context.Context, req *ssov1.GetUserByIdRequest) (*ssov1.GetUserByIdResponse, error) {
	user, err := s.UserInfo.GetUserById(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	respUser := &ssov1.User{
		Id:        user.Id,
		Name:      user.Name,
		FullName:  user.FullName,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		AvatarUrl: user.AvatarUrl,
	}

	return &ssov1.GetUserByIdResponse{User: respUser}, nil
}

func (s *UserInfoAPI) GetUserByEmail(ctx context.Context, req *ssov1.GetUserByEmailRequest) (*ssov1.GetUserByEmailResponse, error) {
	user, err := s.UserInfo.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	respUser := &ssov1.User{
		Id:        user.Id,
		Name:      user.Name,
		FullName:  user.FullName,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		AvatarUrl: user.AvatarUrl,
	}

	return &ssov1.GetUserByEmailResponse{User: respUser}, nil
}

func (s *UserInfoAPI) UpdateUser(ctx context.Context, req *ssov1.UpdateUserRequest) (*ssov1.UpdateUserResponse, error) {
	uptades := req.GetUpdates()

	if uptades == nil {
		return nil, status.Error(codes.InvalidArgument, "updates are empty")
	}

	// TODO: validate

	user, err := s.UserInfo.UpdateUser(ctx, req.GetUserId(), uptades)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNameIsTaken):
			return nil, status.Error(codes.AlreadyExists, "name is taken")
		case errors.Is(err, models.ErrInvalidArgument):
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		case errors.Is(err, models.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	respUser := &ssov1.User{
		Id:        user.Id,
		Name:      user.Name,
		FullName:  user.FullName,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		AvatarUrl: user.AvatarUrl,
	}

	return &ssov1.UpdateUserResponse{User: respUser}, nil

}

func (s *UserInfoAPI) DeleteUser(ctx context.Context, req *ssov1.DeleteUserRequest) (*ssov1.DeleteUserResponse, error) {
	err := s.UserInfo.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return &ssov1.DeleteUserResponse{Success: false}, status.Errorf(codes.NotFound, "user not found")
		}
		return &ssov1.DeleteUserResponse{Success: false}, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.DeleteUserResponse{Success: true}, nil
}
