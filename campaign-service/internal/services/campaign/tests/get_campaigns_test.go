package tests

import (
	"campaigntool/internal/domain/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCreatedCampaign_Success(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		CreatedCampaigns(ctx, userId).
		Return([]*models.Campaign{
			{
				Id: campaignId,
			},
		}, nil)

	campaigns, err := service.GetCreatedCampaigns(ctx, userId)
	require.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, campaignId, campaigns[0].Id)
}

func TestGetCreatedCampaign_NoCampaigns(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()

	userId := int64(123)

	mockGameProvider.EXPECT().
		CreatedCampaigns(ctx, userId).
		Return(nil, models.ErrNoCampaigns)

	campaigns, err := service.GetCreatedCampaigns(ctx, userId)
	require.Error(t, err)
	assert.Equal(t, models.ErrNoCampaigns, errors.Unwrap(err))
	assert.Nil(t, campaigns)
}

func TestGetCurrentCampaign_Success(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		CurrentCampaigns(ctx, userId).
		Return([]*models.CampaignForPlayer{
			{
				Id: campaignId,
			},
		}, nil)

	campaigns, err := service.GetCurrentCampaigns(ctx, userId)
	require.NoError(t, err)
	assert.Equal(t, 1, len(campaigns))
	assert.Equal(t, campaignId, campaigns[0].Id)
}

func TestGetCurrentCampaign_NoCampaigns(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()

	userId := int64(123)

	mockGameProvider.EXPECT().
		CurrentCampaigns(ctx, userId).
		Return(nil, models.ErrNoCampaigns)

	campaigns, err := service.GetCurrentCampaigns(ctx, userId)
	require.Error(t, err)
	assert.Equal(t, models.ErrNoCampaigns, errors.Unwrap(err))
	assert.Nil(t, campaigns)
}
