package tests

import (
	"campaigntool/internal/services/campaign"
	"campaigntool/tests/mocks"
	"log/slog"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func setupTest(t *testing.T) (*campaign.CampaignTool, *mocks.MockGameSaver, *mocks.MockGameProvider) {
	ctrl := gomock.NewController(t)
	mockGameSaver := mocks.NewMockGameSaver(ctrl)
	mockGameProvider := mocks.NewMockGameProvider(ctrl)

	service := &campaign.CampaignTool{
		GameSaver:    mockGameSaver,
		GameProvider: mockGameProvider,
		Log:          slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	return service, mockGameSaver, mockGameProvider
}