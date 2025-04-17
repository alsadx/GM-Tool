package tests

import (
	"campaigntool/internal/domain/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCampaign_Success(t *testing.T) {
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
		DeleteCampaign(ctx, campaignId, 1).
		Return(nil)

	err = service.DeleteCampaign(ctx, campaignId, 1)

	require.NoError(t, err)
}

func TestDeleteCampaign_CampaignNotFound(t *testing.T) {
	service, mockGameSaver, _ := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)
	campaignId := int32(123)

	mockGameSaver.EXPECT().
		DeleteCampaign(ctx, campaignId, 1).
		Return(models.ErrCampaignNotFound)

	err := service.DeleteCampaign(ctx, campaignId, 1)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}

func TestDeleteCampaign_CampaignNotFound2(t *testing.T) {
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
		DeleteCampaign(ctx, campaignId, 2).
		Return(models.ErrCampaignNotFound)

	err = service.DeleteCampaign(ctx, campaignId, 2)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}
