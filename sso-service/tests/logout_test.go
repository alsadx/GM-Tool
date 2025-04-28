package tests

import (
	"sso/tests/suite"
	"testing"

	"protos/gen/go/ssov1"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogout_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()
	name := gofakeit.Name()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})

	require.NoError(t, err)
	userId := respReg.GetUserId()
	require.NotEmpty(t, userId)

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respLog.GetToken())
	assert.NotEmpty(t, respLog.GetRefreshToken())

	respLogout, err := st.AuthClient.Logout(ctx, &ssov1.LogoutRequest{
		UserId: userId,
	})

	require.NoError(t, err)
	assert.Equal(t, true, respLogout.GetSuccess())
}

func TestLogout_UserNotFound(t *testing.T) {
	ctx, st := suite.New(t)

	respLogout, err := st.AuthClient.Logout(ctx, &ssov1.LogoutRequest{
		UserId: int64(0),
	})

	require.Error(t, err)
	assert.False(t, respLogout.GetSuccess())

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, stt.Code(), "unexpected error code")
	require.Equal(t, "user not found", stt.Message(), "unexpected error message")
}
