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

	ctx := context.Background()
	name := "Test Campaign"
	description := "This is a test campaign"
	expectedCampaignId := int64(123)
	userId := int64(123)

	mockGameSaver.EXPECT().
		SaveCampaign(ctx, name, description, userId).
		Return(expectedCampaignId, nil)

	campaignId, err := service.CreateCampaign(ctx, name, description, userId)

	require.NoError(t, err)
	assert.Equal(t, expectedCampaignId, campaignId)
}

func TestCreateCampaign_CampaignAlreadyExists(t *testing.T) {
	service, mockGameSaver, _ := setupTest(t)

	ctx := context.Background()
	name := "Test Campaign"
	description := "This is a test campaign"
	expectedCampaignId := int64(123)
	userId := int64(123)

	mockGameSaver.EXPECT().
		SaveCampaign(ctx, name, description, userId).
		Return(expectedCampaignId, nil)

	campaignId, err := service.CreateCampaign(ctx, name, description, userId)

	require.NoError(t, err)
	assert.Equal(t, expectedCampaignId, campaignId)

	mockGameSaver.EXPECT().
		SaveCampaign(ctx, name, description, userId).
		Return(int64(0), models.ErrCampaignExists)

	campaignId, err = service.CreateCampaign(ctx, name, description, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignExists, errors.Unwrap(err))
	assert.Equal(t, int64(0), campaignId)
}
