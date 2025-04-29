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

	"sso/protos/ssov1"

	"github.com/brianvoe/gofakeit"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestGRPC_GetUser_Success(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserProvider.EXPECT().
		UserByEmail(gomock.Any(), user.Email).
		Return(user, nil)

	mockUserProvider.EXPECT().
		UserById(gomock.Any(), user.Id).
		Return(user, nil)

	respEmail, err := client.GetUserByEmail(ctx, &ssov1.GetUserByEmailRequest{
		Email: user.Email,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respEmail)
	respUser := respEmail.GetUser()
	assert.Equal(t, respUser.Id, user.Id)
	assert.Equal(t, respUser.Name, user.Name)
	assert.Equal(t, respUser.FullName, user.FullName)
	assert.Equal(t, respUser.Email, user.Email)
	assert.Equal(t, respUser.IsAdmin, user.IsAdmin)
	assert.Equal(t, respUser.AvatarUrl, user.AvatarUrl)

	respId, err := client.GetUserById(ctx, &ssov1.GetUserByIdRequest{
		UserId: user.Id,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respId)
	respUser = respId.GetUser()
	assert.Equal(t, respUser.Id, user.Id)
	assert.Equal(t, respUser.Name, user.Name)
	assert.Equal(t, respUser.FullName, user.FullName)
	assert.Equal(t, respUser.Email, user.Email)
	assert.Equal(t, respUser.IsAdmin, user.IsAdmin)
	assert.Equal(t, respUser.AvatarUrl, user.AvatarUrl)
}

func TestGRPC_GetUser_UserNotFound(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockUserProvider.EXPECT().
		UserByEmail(gomock.Any(), user.Email).
		Return(&models.User{}, models.ErrUserNotFound)

	mockUserProvider.EXPECT().
		UserById(gomock.Any(), user.Id).
		Return(&models.User{}, models.ErrUserNotFound)

	respEmail, err := client.GetUserByEmail(ctx, &ssov1.GetUserByEmailRequest{
		Email: user.Email,
	})

	require.Error(t, err)
	require.Empty(t, respEmail)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, models.ErrUserNotFound.Error(), st.Message(), "unexpected error message")

	respId, err := client.GetUserById(ctx, &ssov1.GetUserByIdRequest{
		UserId: user.Id,
	})

	require.Error(t, err)
	require.Empty(t, respId)

	st, ok = status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, models.ErrUserNotFound.Error(), st.Message(), "unexpected error message")
}
