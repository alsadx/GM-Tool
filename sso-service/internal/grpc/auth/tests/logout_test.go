package tests

import (
	"context"
	"net"
	"os"
	"sso/internal/domain/models"
	grpcauth "sso/internal/grpc/auth"
	"sso/tests/setup"
	"testing"
	"time"

	"sso/protos/ssov1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestGRPC_Logout_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, _, _, _ := setup.TestAuth(t)
	srv := grpcauth.ServerAPI{
		Auth: service,
	}

	ssov1.RegisterAuthServer(server, &srv)

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

	client := ssov1.NewAuthClient(clientConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserSaver.EXPECT().
		DeleteSession(gomock.Any(), int64(1)).
		Return(nil)

	resp, err := client.Logout(ctx, &ssov1.LogoutRequest{
		UserId: int64(1),
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestGRPC_Logout_UserNotFound(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, _, _, _ := setup.TestAuth(t)
	srv := grpcauth.ServerAPI{
		Auth: service,
	}

	ssov1.RegisterAuthServer(server, &srv)

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

	client := ssov1.NewAuthClient(clientConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserSaver.EXPECT().
		DeleteSession(gomock.Any(), int64(1)).
		Return(models.ErrUserNotFound)

	resp, err := client.Logout(ctx, &ssov1.LogoutRequest{
		UserId: int64(1),
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code: %v", st.Code())
	require.Equal(t, "user not found", st.Message(), "unexpected error message")
}
