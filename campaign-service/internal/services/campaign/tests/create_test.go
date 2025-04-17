package tests

import (
	"campaigntool/internal/domain/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCampaign_Success(t *testing.T) {
	service, mockGameSaver, _ := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)
	name := "Test Campaign"
	description := "This is a test campaign"
	expectedCampaignId := int32(123)

	mockGameSaver.EXPECT().
		SaveCampaign(ctx, name, description, 1).
		Return(expectedCampaignId, nil)

	campaignId, err := service.CreateCampaign(ctx, name, description, 1)

	require.NoError(t, err)
	assert.Equal(t, expectedCampaignId, campaignId)
}

func TestCreateCampaign_CampaignAlreadyExists(t *testing.T) {
	service, mockGameSaver, _ := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)
	name := "Test Campaign"
	description := "This is a test campaign"
	expectedCampaignId := int32(123)

	mockGameSaver.EXPECT().
		SaveCampaign(ctx, name, description, 1).
		Return(expectedCampaignId, nil)

	campaignId, err := service.CreateCampaign(ctx, name, description, 1)

	require.NoError(t, err)
	assert.Equal(t, expectedCampaignId, campaignId)

	mockGameSaver.EXPECT().
		SaveCampaign(ctx, name, description, 1).
		Return(int32(0), models.ErrCampaignExists)

	campaignId, err = service.CreateCampaign(ctx, name, description, 1)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignExists, errors.Unwrap(err))
	assert.Equal(t, int32(0), campaignId)
}
