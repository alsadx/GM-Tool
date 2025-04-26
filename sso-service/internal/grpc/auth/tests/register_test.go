package tests

import (
	"auth"
	"context"
	"net"
	"os"
	"sso/tests/setup"
	"testing"
	"time"
	grpcauth"sso/internal/grpc/auth"
	ssov1 "github.com/alsadx/protos/gen/go/sso"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestGRPC_Register_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(auth.AuthInterceptor))
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

	name := "name"
	email := "email"
	password := "password"
	expectedUserId := int32(123)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	mockHasher.EXPECT().
		Hash(password).
		Return([]byte("hashed-password"), nil)

	mockUserSaver.EXPECT().
		SaveUser(gomock.Any(), email, name, password).
		Return(expectedUserId, nil)

	resp, err := client.Register(ctx, &ssov1.RegisterRequest{
		Name:        name,
		Email:       email,
		Password:    password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetUserId())
}
