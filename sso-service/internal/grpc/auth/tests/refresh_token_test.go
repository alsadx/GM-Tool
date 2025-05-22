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

func TestGRPC_RefreshToken_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, mockUserProvider, _, mockTokenManager := setup.TestAuth(t)
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

	refreshToken := "refresh_token"
	user := &models.User{
		Id:        1,
		Email:     "email",
		PassHash:  []byte("password"),
		Name:      "name",
		FullName:  "full_name",
		IsAdmin:   false,
		AvatarUrl: "avatar",
	}

	mockUserProvider.EXPECT().
		UserByRefreshToken(gomock.Any(), refreshToken).
		Return(user, nil)

	mockTokenManager.EXPECT().
		NewJWT(user, "secret", gomock.Any()).
		Return("access_token", nil)

	mockTokenManager.EXPECT().
		NewRefreshToken().
		Return("new_refresh_token", nil)

	mockUserSaver.EXPECT().
		SetSession(gomock.Any(), user.Id, gomock.Any()).
		Return(nil)

	resp, err := client.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})

	require.NoError(t, err)
	assert.Equal(t, "access_token", resp.GetToken())
	assert.Equal(t, "new_refresh_token", resp.GetRefreshToken())
}

func TestGRPC_RefreshToken_UserNotFound(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, _, mockUserProvider, _, _ := setup.TestAuth(t)
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

	refreshToken := "refresh_token"

	mockUserProvider.EXPECT().
		UserByRefreshToken(gomock.Any(), refreshToken).
		Return(&models.User{}, models.ErrUserNotFound)

	resp, err := client.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})

	require.Error(t, err)
	assert.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.Unauthenticated, st.Code(), "unexpected error code")
	require.Equal(t, "invalid or expired refresh token", st.Message(), "unexpected error message")
}
