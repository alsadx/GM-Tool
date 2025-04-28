package setup

import (
	"log/slog"
	"os"
	userinfo "sso/internal/services/user_info"
	"sso/internal/services/auth"
	"testing"

	"sso/tests/mocks"

	"github.com/golang/mock/gomock"
)

func TestAuth(t *testing.T) (*auth.Auth, *mocks.MockUserSaver, *mocks.MockUserProvider, *mocks.MockHasher, *mocks.MockTokenManager) {
	ctrl := gomock.NewController(t)
	mockUserSaver := mocks.NewMockUserSaver(ctrl)
	mockUserProvider := mocks.NewMockUserProvider(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)
	mockTokenManager := mocks.NewMockTokenManager(ctrl)

	service := &auth.Auth{
		UserSaver:    mockUserSaver,
		UserProvider: mockUserProvider,
		Log:          slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Hasher:       mockHasher,
		TokenManager: mockTokenManager,
	}

	return service, mockUserSaver, mockUserProvider, mockHasher, mockTokenManager
}

func TestUser(t *testing.T) (*userinfo.UserInfo, *mocks.MockUserSaver, *mocks.MockUserProvider) {
	ctrl := gomock.NewController(t)
	mockUserSaver := mocks.NewMockUserSaver(ctrl)
	mockUserProvider := mocks.NewMockUserProvider(ctrl)

	service := &userinfo.UserInfo{
		UserSaver:    mockUserSaver,
		UserProvider: mockUserProvider,
		Log:          slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	return service, mockUserSaver, mockUserProvider
}
