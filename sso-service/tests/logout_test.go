package tests

import (
	"sso/internal/lib/jwt"
	"sso/tests/suite"
	"testing"
	"time"

	"protos/gen/go/ssov1"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestRegisterLogin_Logout_HappyPath(t *testing.T) {
	ctx, suite := suite.New(t)
	tokenManager := jwt.NewTokenManager()

	email := gofakeit.Email()
	password := randomFakePassword()
	name := gofakeit.Name()

	respReg, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLog, err := suite.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLog.GetToken()
	require.NotEmpty(t, token)

	parsedJWT, err := tokenManager.ParseJWT(respLog.Token, signKey)
	require.NoError(t, err)

	assert.Equal(t, respReg.GetUserId(), parsedJWT.UserId)
	assert.Equal(t, email, parsedJWT.Email)

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(suite.Cfg.Auth.AccessTokenTTL).Unix(), parsedJWT.ExpiresAt.Unix(), deltaSeconds)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	respOut, err := suite.AuthClient.Logout(ctx, &ssov1.LogoutRequest{})
	require.NoError(t, err)

	assert.True(t, respOut.GetSuccess())
}
