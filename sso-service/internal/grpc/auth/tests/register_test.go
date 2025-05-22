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

func TestGRPC_Register_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, _, mockHasher, _ := setup.TestAuth(t)
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
	hashedPassword := []byte("hashed-password")
	expectedUserId := int64(123)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	mockHasher.EXPECT().
		Hash(password).
		Return(hashedPassword, nil)

	mockUserSaver.EXPECT().
		SaveUser(gomock.Any(), email, name, hashedPassword).
		Return(expectedUserId, nil)

	resp, err := client.Register(ctx, &ssov1.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetUserId())
}

func TestGRPC_Register_EmptyName(t *testing.T) {
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

	name := ""
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	resp, err := client.Register(ctx, &ssov1.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Equal(t, int64(0), resp.GetUserId())
	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "name is required", st.Message(), "unexpected error message")
}

func TestGRPC_Register_EmptyEmail(t *testing.T) {
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

	name := gofakeit.Name()
	email := ""
	password := gofakeit.Password(true, true, true, true, false, 6)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	resp, err := client.Register(ctx, &ssov1.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Equal(t, int64(0), resp.GetUserId())
	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "email is required", st.Message(), "unexpected error message")
}

func TestGRPC_Register_WrongEmail(t *testing.T) {
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

	name := gofakeit.Name()
	email := "invalid-email"
	password := gofakeit.Password(true, true, true, true, false, 6)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	resp, err := client.Register(ctx, &ssov1.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Equal(t, int64(0), resp.GetUserId())
	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "email is required", st.Message(), "unexpected error message")
}

func TestGRPC_Register_EmptyPassword(t *testing.T) {
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

	name := gofakeit.Name()
	email := gofakeit.Email()
	password := ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	resp, err := client.Register(ctx, &ssov1.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Equal(t, int64(0), resp.GetUserId())
	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "password is required", st.Message(), "unexpected error message")
}

func TestGRPC_Register_UserExists(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockUserSaver, _, mockHasher, _ := setup.TestAuth(t)
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
	hashedPassword := []byte("hashed-password")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	mockHasher.EXPECT().
		Hash(password).
		Return(hashedPassword, nil)

	mockUserSaver.EXPECT().
		SaveUser(gomock.Any(), email, name, hashedPassword).
		Return(int64(0), models.ErrUserExists)

	resp, err := client.Register(ctx, &ssov1.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Equal(t, int64(0), resp.GetUserId())
	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.AlreadyExists, st.Code(), "unexpected error code")
	require.Equal(t, "user already exists", st.Message(), "unexpected error message")
}