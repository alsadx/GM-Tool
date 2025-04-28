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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestGRPC_DeleteUser_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, _ := setup.TestUser(t)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId := int64(1)

	mockUserSaver.EXPECT().
		DeleteSession(gomock.Any(), userId).
		Return(nil)

	mockUserSaver.EXPECT().
		DeleteUser(gomock.Any(), userId).
		Return(nil)

	resp, err := client.DeleteUser(ctx, &ssov1.DeleteUserRequest{
		UserId: userId,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.True(t, resp.GetSuccess())
}

func TestGRPC_DeleteUser_UserNotFound(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, _ := setup.TestUser(t)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId := int64(1)

	mockUserSaver.EXPECT().
		DeleteSession(gomock.Any(), userId).
		Return(models.ErrUserNotFound)

	resp, err := client.DeleteUser(ctx, &ssov1.DeleteUserRequest{
		UserId: userId,
	})

	require.Error(t, err)
	require.Empty(t, resp)
	require.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code())
	require.Equal(t, "user not found", st.Message())
}

func TestGRPC_DeleteUser_UserNotFound2(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, _ := setup.TestUser(t)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId := int64(1)

	mockUserSaver.EXPECT().
		DeleteSession(gomock.Any(), userId).
		Return(nil)

	mockUserSaver.EXPECT().
		DeleteUser(gomock.Any(), userId).
		Return(models.ErrUserNotFound)

	resp, err := client.DeleteUser(ctx, &ssov1.DeleteUserRequest{
		UserId: userId,
	})

	require.Error(t, err)
	require.Empty(t, resp)
	require.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code())
	require.Equal(t, "user not found", st.Message())
}
