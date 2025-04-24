package tests

import (
	"log/slog"
	"os"
	"sso/internal/services/auth"
	"testing"

	"sso/tests/mocks"

	"github.com/golang/mock/gomock"
)

func setupTest(t *testing.T) (*auth.Auth, *mocks.MockUserSaver, *mocks.MockUserProvider, *mocks.MockHasher) {
	ctrl := gomock.NewController(t)
	mockUserSaver := mocks.NewMockUserSaver(ctrl)
	mockUserProvider := mocks.NewMockUserProvider(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)

	service := &auth.Auth{
		UserSaver:    mockUserSaver,
		UserProvider: mockUserProvider,
		Log:          slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Hasher: mockHasher,
	}

	return service, mockUserSaver, mockUserProvider, mockHasher
}
