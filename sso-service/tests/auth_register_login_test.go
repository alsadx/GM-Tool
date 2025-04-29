package tests

import (
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

const (
	signKey = "secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
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
}

func TestReqister_UserExists(t *testing.T) {
	ctx, suite := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()
	name := gofakeit.Name()

	resp, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetUserId())

	resp, err = suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})
	require.Error(t, err)
	assert.Empty(t, resp.GetUserId())
	
	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.AlreadyExists, stt.Code(), "unexpected error code")
	require.Equal(t, "user already exists", stt.Message(), "unexpected error message")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		username    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			username:    gofakeit.Name(),
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			username:    gofakeit.Name(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with Empty Username",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			username:    "",
			expectedErr: "name is required",
		},
		{
			name:        "Register with All Empty",
			email:       "",
			password:    "",
			username:    "",
			expectedErr: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
				Name:     tt.username,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			expectedErr: "invalid email or password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomFakePassword(),
				Name:     gofakeit.Name(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
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
		Password: "invalid",
	})

	require.Error(t, err)
	assert.Empty(t, respLog.GetToken())

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, stt.Code(), "unexpected error code")
	require.Equal(t, "invalid email or password", stt.Message(), "unexpected error message")
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
