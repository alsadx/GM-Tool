package userinfo

import (
	"context"
	"errors"
	"sso/internal/domain/models"

	ssov1 "github.com/alsadx/protos/gen/go/sso"

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

type userInfoAPI struct {
	ssov1.UnimplementedUserInfoServer
	userInfo UserInfo
}

func RegisterUserInfoAPI(gRPC *grpc.Server, userInfo UserInfo) {
	ssov1.RegisterUserInfoServer(gRPC, &userInfoAPI{userInfo: userInfo})
}

func (s *userInfoAPI) GetUserById(ctx context.Context, req *ssov1.GetUserByIdRequest) (*ssov1.GetUserByIdResponse, error) {
	user, err := s.userInfo.GetUserById(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

    respUser := &ssov1.User{
        Id: user.Id,
        Name:   user.Name,
        FullName: user.Name,
        Email:  user.Email,
        IsAdmin: user.IsAdmin,
        AvatarUrl: user.AvatarUrl,
    }

    return &ssov1.GetUserByIdResponse{User: respUser}, nil
}

func (s *userInfoAPI) GetUserByEmail(ctx context.Context, req *ssov1.GetUserByEmailRequest) (*ssov1.GetUserByEmailResponse, error) {
    user, err := s.userInfo.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

    respUser := &ssov1.User{
        Id: user.Id,
        Name:   user.Name,
        FullName: user.Name,
        Email:  user.Email,
        IsAdmin: user.IsAdmin,
        AvatarUrl: user.AvatarUrl,
    }

    return &ssov1.GetUserByEmailResponse{User: respUser}, nil
}

func (s *userInfoAPI) UpdateUser(ctx context.Context, req *ssov1.UpdateUserRequest) (*ssov1.UpdateUserResponse, error) {
    uptades := req.GetUpdates()

    if uptades == nil {
        return nil, status.Error(codes.InvalidArgument, "updates is required")
    }

    // TODO: validate

    user, err := s.userInfo.UpdateUser(ctx, req.GetUserId(), uptades)
    if err != nil {
        if errors.Is(err, models.ErrUserNotFound) {
            return nil, status.Errorf(codes.NotFound, "user not found")
        } else if errors.Is(err, models.ErrInvalidArgument) {
            return nil, status.Error(codes.InvalidArgument, err.Error())
        }
        return nil, status.Error(codes.Internal, "internal error")
    }

    respUser := &ssov1.User{
        Id: user.Id,
        Name:   user.Name,
        FullName: user.Name,
        Email:  user.Email,
        IsAdmin: user.IsAdmin,
        AvatarUrl: user.AvatarUrl,
    }

    return &ssov1.UpdateUserResponse{User: respUser}, nil

}

func (s *userInfoAPI) DeleteUser(ctx context.Context, req *ssov1.DeleteUserRequest) (*ssov1.DeleteUserResponse, error) {
    err := s.userInfo.DeleteUser(ctx, req.GetUserId())
    if err != nil {
        if errors.Is(err, models.ErrUserNotFound) {
            return &ssov1.DeleteUserResponse{Success: false}, status.Errorf(codes.NotFound, "user not found")
        } 
        return &ssov1.DeleteUserResponse{Success: false}, status.Error(codes.Internal, "internal error")
    }
        
    return &ssov1.DeleteUserResponse{Success: true}, nil
}
