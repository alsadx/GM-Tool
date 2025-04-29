package tests

import (
	"fmt"
	"sso/internal/lib/jwt"
	"sso/tests/suite"
	"testing"
	"time"

	"sso/protos/ssov1"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRefreshToken_HappyPath(t *testing.T) {
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
	assert.NotEmpty(t, respReg.GetUserId())

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respLog.GetToken())
	assert.NotEmpty(t, respLog.GetRefreshToken())

	refreshTime := time.Now()

	respRefresh, err := st.AuthClient.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: respLog.GetRefreshToken(),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respRefresh.GetToken())
	assert.NotEmpty(t, respRefresh.GetRefreshToken())

	assert.NotEqual(t, respLog.GetRefreshToken(), respRefresh.GetRefreshToken())

	tokenManager := jwt.NewTokenManager()

	parsedJWT, err := tokenManager.ParseJWT(respRefresh.Token, signKey)

	require.NoError(t, err)
	assert.Equal(t, respReg.GetUserId(), parsedJWT.UserId)
	assert.Equal(t, email, parsedJWT.Email)

	const deltaSeconds = 1

	assert.InDelta(t, refreshTime.Add(st.Cfg.Auth.AccessTokenTTL).Unix(), parsedJWT.ExpiresAt.Unix(), deltaSeconds)
}

func TestRefreshToken_InvalidRefreshToken(t *testing.T) {
	ctx, st := suite.New(t)

	respRefresh, err := st.AuthClient.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: "invalid_refresh_token",
	})

	require.Error(t, err)
	assert.Empty(t, respRefresh)

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.Unauthenticated, stt.Code(), "unexpected error code")
	require.Equal(t, "invalid or expired refresh token", stt.Message(), "unexpected error message")
}

func TestRefreshToken_ExpiredRefreshToken(t *testing.T) {
	// Set refresh token TTL to 5 seconds!!

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
	assert.NotEmpty(t, respReg.GetUserId())

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respLog.GetToken())
	assert.NotEmpty(t, respLog.GetRefreshToken())

	fmt.Println("waiting for refresh token to expire...")
	time.Sleep(5 * time.Second)

	respRefresh, err := st.AuthClient.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: respLog.GetRefreshToken(),
	})

	require.Error(t, err)
	assert.Empty(t, respRefresh)

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.Unauthenticated, stt.Code(), "unexpected error code")
	require.Equal(t, "invalid or expired refresh token", stt.Message(), "unexpected error message")
}
