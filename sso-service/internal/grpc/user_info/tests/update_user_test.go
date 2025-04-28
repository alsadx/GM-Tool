package tests

import (
	"context"
	"net"
	"os"
	"sso/internal/domain/models"
	grpcuserinfo "sso/internal/grpc/user_info"
	"sso/tests/setup"
	"testing"
	"time"

	"protos/gen/go/ssov1"

	"github.com/brianvoe/gofakeit"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestGRPC_UpdateUser_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, mockUserProvider := setup.TestUser(t)
	srv := grpcuserinfo.UserInfoAPI{
		UserInfo: service,
	}

	ssov1.RegisterUserInfoServer(server, &srv)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(listener)
	defer server.Stop()

	serverAddress := listener.Addr().String()

	clientConn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	require.NoError(t, err)
	defer clientConn.Close()

	client := ssov1.NewUserInfoClient(clientConn)

	user := &models.User{
		Id:        1,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		PassHash:  []byte(gofakeit.Password(true, true, true, true, false, 6)),
		FullName:  "",
		IsAdmin:   false,
		AvatarUrl: "",
	}

	updates := map[string]string{
		"name":      user.Name+"New",
		"full_name": user.FullName+"New",
		"avatar_url":    user.AvatarUrl+"New",
	}

	updatedUser := &models.User{
		Id:        user.Id,
		Name:      user.Name+"New",
		Email:     user.Email,
		PassHash:  user.PassHash,
		FullName:  user.FullName+"New",
		IsAdmin:   user.IsAdmin,
		AvatarUrl: user.AvatarUrl+"New",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserProvider.EXPECT().
		UserById(gomock.Any(), user.Id).
		Return(user, nil)

	mockUserSaver.EXPECT().
		UpdateUser(gomock.Any(), user).
		Return(nil)

	resp, err := client.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: user.Id,
		Updates: updates,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp)

	respUser := resp.GetUser()
	assert.Equal(t, user.Id, respUser.Id)
	assert.Equal(t, updatedUser.Name, respUser.Name)
	assert.Equal(t, updatedUser.FullName, respUser.FullName)
	assert.Equal(t, updatedUser.AvatarUrl, respUser.AvatarUrl)
	assert.Equal(t, updatedUser.Email, respUser.Email)
	assert.Equal(t, updatedUser.IsAdmin, respUser.IsAdmin)
}

func TestGRPC_UpdateUser_UserNotFound(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, _, mockUserProvider := setup.TestUser(t)
	srv := grpcuserinfo.UserInfoAPI{
		UserInfo: service,
	}

	ssov1.RegisterUserInfoServer(server, &srv)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(listener)
	defer server.Stop()

	serverAddress := listener.Addr().String()

	clientConn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	require.NoError(t, err)
	defer clientConn.Close()

	client := ssov1.NewUserInfoClient(clientConn)

	user := &models.User{
		Id:        1,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		PassHash:  []byte(gofakeit.Password(true, true, true, true, false, 6)),
		FullName:  "",
		IsAdmin:   false,
		AvatarUrl: "",
	}

	updates := map[string]string{
		"name":      user.Name+"New",
		"full_name": user.FullName+"New",
		"avatar_url":    user.AvatarUrl+"New",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserProvider.EXPECT().
		UserById(gomock.Any(), user.Id).
		Return(&models.User{}, models.ErrUserNotFound)

	resp, err := client.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: user.Id,
		Updates: updates,
	})

	require.Error(t, err)
	require.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "user not found", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateUser_InvalidArguments(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, mockUserProvider := setup.TestUser(t)
	srv := grpcuserinfo.UserInfoAPI{
		UserInfo: service,
	}

	ssov1.RegisterUserInfoServer(server, &srv)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(listener)
	defer server.Stop()

	serverAddress := listener.Addr().String()

	clientConn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	require.NoError(t, err)
	defer clientConn.Close()

	client := ssov1.NewUserInfoClient(clientConn)

	user := &models.User{
		Id:        1,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		PassHash:  []byte(gofakeit.Password(true, true, true, true, false, 6)),
		FullName:  "",
		IsAdmin:   false,
		AvatarUrl: "",
	}

	updates := map[string]string{
		"name":      user.Name+"New",
		"full_name": user.FullName+"New",
		"avatar_url":    user.AvatarUrl+"New",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserProvider.EXPECT().
		UserById(gomock.Any(), user.Id).
		Return(user, nil)

	mockUserSaver.EXPECT().
		UpdateUser(gomock.Any(), user).
		Return(models.ErrInvalidArgument)

	resp, err := client.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: user.Id,
		Updates: updates,
	})

	require.Error(t, err)
	require.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "invalid argument", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateUser_EmptyUpdates(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, _, _ := setup.TestUser(t)
	srv := grpcuserinfo.UserInfoAPI{
		UserInfo: service,
	}

	ssov1.RegisterUserInfoServer(server, &srv)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(listener)
	defer server.Stop()

	serverAddress := listener.Addr().String()

	clientConn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	require.NoError(t, err)
	defer clientConn.Close()

	client := ssov1.NewUserInfoClient(clientConn)

	userId := int64(1)

	updates := map[string]string{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: userId,
		Updates: updates,
	})

	require.Error(t, err)
	require.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "updates are empty", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateUser_NoUpdates(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, _, mockUserProvider := setup.TestUser(t)
	srv := grpcuserinfo.UserInfoAPI{
		UserInfo: service,
	}

	ssov1.RegisterUserInfoServer(server, &srv)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(listener)
	defer server.Stop()

	serverAddress := listener.Addr().String()

	clientConn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	require.NoError(t, err)
	defer clientConn.Close()

	client := ssov1.NewUserInfoClient(clientConn)

	user := &models.User{
		Id:        1,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		PassHash:  []byte(gofakeit.Password(true, true, true, true, false, 6)),
		FullName:  "",
		IsAdmin:   false,
		AvatarUrl: "",
	}

	updates := map[string]string{
		"name":      user.Name,
		"full_name": user.FullName,
		"avatar_url":    user.AvatarUrl,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserProvider.EXPECT().
		UserById(gomock.Any(), user.Id).
		Return(user, nil)

	resp, err := client.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: user.Id,
		Updates: updates,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp)

	respUser := resp.GetUser()
	assert.Equal(t, user.Id, respUser.Id)
	assert.Equal(t, user.Name, respUser.Name)
	assert.Equal(t, user.FullName, respUser.FullName)
	assert.Equal(t, user.AvatarUrl, respUser.AvatarUrl)
	assert.Equal(t, user.Email, respUser.Email)
	assert.Equal(t, user.IsAdmin, respUser.IsAdmin)
}

func TestGRPC_UpdateUser_NameIsTaken(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, mockUserProvider := setup.TestUser(t)
	srv := grpcuserinfo.UserInfoAPI{
		UserInfo: service,
	}

	ssov1.RegisterUserInfoServer(server, &srv)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(listener)
	defer server.Stop()

	serverAddress := listener.Addr().String()

	clientConn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	require.NoError(t, err)
	defer clientConn.Close()

	client := ssov1.NewUserInfoClient(clientConn)

	user := &models.User{
		Id:        1,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		PassHash:  []byte(gofakeit.Password(true, true, true, true, false, 6)),
		FullName:  "",
		IsAdmin:   false,
		AvatarUrl: "",
	}

	updates := map[string]string{
		"name": gofakeit.Name(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserProvider.EXPECT().
		UserById(gomock.Any(), user.Id).
		Return(user, nil)

	mockUserSaver.EXPECT().
		UpdateUser(gomock.Any(), user).
		Return(models.ErrNameIsTaken)

	resp, err := client.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: user.Id,
		Updates: updates,
	})

	require.Error(t, err)
	require.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.AlreadyExists, st.Code(), "unexpected error code")
	require.Equal(t, "name is taken", st.Message(), "unexpected error message")
}