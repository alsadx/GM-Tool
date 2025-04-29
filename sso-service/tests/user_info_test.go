package tests

import (
	"sso/tests/suite"
	"strings"
	"testing"

	"sso/protos/ssov1"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetUser_UpdateUser_HappyPath(t *testing.T) {
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
	assert.NotEmpty(t, respReg.UserId)

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respLog.GetToken())
	assert.NotEmpty(t, respLog.GetRefreshToken())

	respGetById, err := st.UserInfoClient.GetUserById(ctx, &ssov1.GetUserByIdRequest{
		UserId: userId,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respGetById)

	respUser := respGetById.GetUser()
	assert.Equal(t, userId, respUser.Id)
	assert.Equal(t, email, respUser.Email)
	assert.Equal(t, name, respUser.Name)
	assert.Equal(t, "", respUser.FullName)
	assert.Equal(t, "", respUser.AvatarUrl)
	assert.False(t, respUser.IsAdmin)

	respGetByEmail, err := st.UserInfoClient.GetUserByEmail(ctx, &ssov1.GetUserByEmailRequest{
		Email: email,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respGetById)

	respUser = respGetByEmail.GetUser()
	assert.Equal(t, userId, respUser.Id)
	assert.Equal(t, email, respUser.Email)
	assert.Equal(t, name, respUser.Name)
	assert.Equal(t, "", respUser.FullName)
	assert.Equal(t, "", respUser.AvatarUrl)
	assert.False(t, respUser.IsAdmin)

	fullName := "New Full Name"
	avatarUrl := "path/to/avatar.jpg"

	updates := map[string]string{
		"full_name":  fullName,
		"avatar_url": avatarUrl,
	}

	respUpdate, err := st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId:  userId,
		Updates: updates,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respUpdate)

	respGetById, err = st.UserInfoClient.GetUserById(ctx, &ssov1.GetUserByIdRequest{
		UserId: userId,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respGetById)

	respUser = respGetById.GetUser()
	assert.Equal(t, userId, respUser.Id)
	assert.Equal(t, email, respUser.Email)
	assert.Equal(t, name, respUser.Name)
	assert.Equal(t, fullName, respUser.FullName)
	assert.Equal(t, avatarUrl, respUser.AvatarUrl)
	assert.False(t, respUser.IsAdmin)

	respGetByEmail, err = st.UserInfoClient.GetUserByEmail(ctx, &ssov1.GetUserByEmailRequest{
		Email: email,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respGetById)

	respUser = respGetByEmail.GetUser()
	assert.Equal(t, userId, respUser.Id)
	assert.Equal(t, email, respUser.Email)
	assert.Equal(t, name, respUser.Name)
	assert.Equal(t, fullName, respUser.FullName)
	assert.Equal(t, avatarUrl, respUser.AvatarUrl)
	assert.False(t, respUser.IsAdmin)
}

func TestGetUser_UserNotFound(t *testing.T) {
	ctx, st := suite.New(t)

	respId, err := st.UserInfoClient.GetUserById(ctx, &ssov1.GetUserByIdRequest{
		UserId: 0,
	})

	require.Error(t, err)
	require.Empty(t, respId)

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, stt.Code(), "unexpected error code")
	require.Equal(t, "user not found", stt.Message(), "unexpected error message")

	respEmail, err := st.UserInfoClient.GetUserByEmail(ctx, &ssov1.GetUserByEmailRequest{
		Email: "invalid_email",
	})

	require.Error(t, err)
	require.Empty(t, respEmail)

	stt, ok = status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.NotFound, stt.Code(), "unexpected error code")
	require.Equal(t, "user not found", stt.Message(), "unexpected error message")
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	ctx, st := suite.New(t)

	updates := map[string]string{
		"full_name":  "New Full Name",
		"avatar_url": "path/to/avatar.jpg",
	}

	resp, err := st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId:  0,
		Updates: updates,
	})

	require.Error(t, err)
	require.Empty(t, resp)

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.NotFound, stt.Code(), "unexpected error code")
	require.Equal(t, "user not found", stt.Message(), "unexpected error message")
}

func TestUpdateUser_EmptyUpdates(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId:  1,
		Updates: map[string]string{},
	})

	require.Error(t, err)
	require.Empty(t, resp)

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, stt.Code(), "unexpected error code")
	require.Equal(t, "updates are empty", stt.Message(), "unexpected error message")
}

func TestUpdateUser_NameIsTaken(t *testing.T) {
	ctx, st := suite.New(t)

	takenName := gofakeit.Name()

	respReq, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    gofakeit.Email(),
		Password: randomFakePassword(),
		Name:     takenName,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReq.GetUserId())

	email := gofakeit.Email()
	password := randomFakePassword()

	respReq, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     gofakeit.Name(),
	})

	require.NoError(t, err)
	userId := respReq.GetUserId()
	assert.NotEmpty(t, userId)

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respLog.GetToken())
	assert.NotEmpty(t, respLog.GetRefreshToken())

	updates := map[string]string{
		"name":       takenName,
		"avatar_url": "path/to/avatar.jpg",
	}

	resp, err := st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId:  userId,
		Updates: updates,
	})

	require.Error(t, err)
	require.Empty(t, resp)

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.AlreadyExists, stt.Code(), "unexpected error code")
	require.Equal(t, "name is taken", stt.Message(), "unexpected error message")
}

func TestUpdateUser_NoUpdates(t *testing.T) {
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
	assert.NotEmpty(t, respReg.UserId)

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respLog.GetToken())
	assert.NotEmpty(t, respLog.GetRefreshToken())

	updates := map[string]string{
		"name": name,
	}

	resp, err := st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId:  userId,
		Updates: updates,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp)
}

func TestUpdateUser_InvalidArgument(t *testing.T) {
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
	assert.NotEmpty(t, respReg.UserId)

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respLog.GetToken())
	assert.NotEmpty(t, respLog.GetRefreshToken())

	resp, err := st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: userId,
		Updates: map[string]string{
			"full_name":  "New Full Name",
			"avatar_url": "",
		},
	})

	require.Error(t, err)
	require.Empty(t, resp)

	stt, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, stt.Code(), "unexpected error code")
	require.Equal(t, "invalid argument", stt.Message(), "unexpected error message")

	resp, err = st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: userId,
		Updates: map[string]string{
			"full_name":  "",
			"avatar_url": "path/to/avatar.jpg",
		},
	})

	require.Error(t, err)
	require.Empty(t, resp)

	stt, ok = status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, stt.Code(), "unexpected error code")
	require.Equal(t, "invalid argument", stt.Message(), "unexpected error message")

	resp, err = st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: userId,
		Updates: map[string]string{
			"full_name": "New Full Name",
			"name":      "",
		},
	})

	require.Error(t, err)
	require.Empty(t, resp)

	stt, ok = status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, stt.Code(), "unexpected error code")
	require.Equal(t, "invalid argument", stt.Message(), "unexpected error message")

	var builder strings.Builder

	for i := 0; i < 256; i++ {
		builder.WriteByte('a')
	}

	longString := builder.String()

	resp, err = st.UserInfoClient.UpdateUser(ctx, &ssov1.UpdateUserRequest{
		UserId: userId,
		Updates: map[string]string{
			"full_name": longString,
		},
	})

	require.Error(t, err)
	require.Empty(t, resp)

	stt, ok = status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, stt.Code(), "unexpected error code")
	require.Equal(t, "invalid argument", stt.Message(), "unexpected error message")
}
