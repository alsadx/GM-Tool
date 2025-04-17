package tests

import (
	"campaigntool/internal/domain/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCreatedCampaign_Success(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)
	campaignId := int32(123)

	mockGameProvider.EXPECT().
		CreatedCampaigns(ctx, 1).
		Return([]models.Campaign{
			{
				Id: campaignId,
			},
		}, nil)

	campaigns, err := service.GetCreatedCampaigns(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, campaignId, campaigns[0].Id)
}

func TestGetCurrentCampaign_Success(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)
	campaignId := int32(123)

	mockGameProvider.EXPECT().
		CurrentCampaigns(ctx, 1).
		Return([]models.CampaignForPlayer{
			{
				Id: campaignId,
			},
		}, nil)

	campaigns, err := service.GetCurrentCampaigns(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, campaignId, campaigns[0].Id)
}