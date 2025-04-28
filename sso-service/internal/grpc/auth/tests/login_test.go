package tests

import (
	"context"
	"errors"
	"net"
	"os"
	"sso/internal/domain/models"
	grpcauth "sso/internal/grpc/auth"
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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestGRPC_Login_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, mockUserProvider, mockHasher, mockTokenManager := setup.TestAuth(t)
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

	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)
	hashedPassword := []byte("$2a$10$vR/gn5MPG9g5JVZPlhj")

	user := &models.User{
		Id:        1,
		Name:      name,
		Email:     email,
		PassHash:  hashedPassword,
		FullName:  "",
		IsAdmin:   false,
		AvatarUrl: "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	mockUserProvider.EXPECT().
		UserByEmail(gomock.Any(), email).
		Return(user, nil)

	mockHasher.EXPECT().
		CheckPassword(password, hashedPassword).
		Return(nil)

	mockTokenManager.EXPECT().
		NewJWT(user, "secret", gomock.Any()).
		Return("token", nil)

	mockTokenManager.EXPECT().
		NewRefreshToken().
		Return("refresh_token", nil)

	mockUserSaver.EXPECT().
		SetSession(gomock.Any(), int64(1), gomock.Any()).
		Return(nil)

	resp, err := client.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.Equal(t, "token", resp.GetToken())
	assert.Equal(t, "refresh_token", resp.GetRefreshToken())
}

func TestGRPC_Login_InvalidPassword(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, _, mockUserProvider, mockHasher, _ := setup.TestAuth(t)
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

	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)
	hashedPassword := []byte("$2a$10$vR/gn5MPG9g5JVZPlhj")

	user := &models.User{
		Id:        1,
		Name:      name,
		Email:     email,
		PassHash:  hashedPassword,
		FullName:  "",
		IsAdmin:   false,
		AvatarUrl: "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	mockUserProvider.EXPECT().
		UserByEmail(gomock.Any(), email).
		Return(user, nil)

	mockHasher.EXPECT().
		CheckPassword(password, hashedPassword).
		Return(errors.New("invalid password"))

	resp, err := client.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "invalid email or password", st.Message(), "unexpected error message")
}

func TestGRPC_Login_EmptyEmail(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, _, _, _, _ := setup.TestAuth(t)
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

	email := ""
	password := gofakeit.Password(true, true, true, true, false, 6)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	resp, err := client.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "email is required", st.Message(), "unexpected error message")
}
