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
	service, mockGameSaver, mockGameProvider := setupTest(t)
	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaignPlayers(ctx, campaignId).
		Return([]int64{}, nil)

	mockGameSaver.EXPECT().
		DeleteCampaign(ctx, campaignId, userId).
		Return(nil)

	err := service.DeleteCampaign(ctx, campaignId, userId)

	require.NoError(t, err)
}

func TestDeleteCampaign_WithPlayers(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)
	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaignPlayers(ctx, campaignId).
		Return([]int64{1, 2}, nil)

	mockGameSaver.EXPECT().
		RemovePlayer(ctx, campaignId, int64(1)).
		Return(nil)

	mockGameSaver.EXPECT().
		RemovePlayer(ctx, campaignId, int64(2)).
		Return(nil)

	mockGameSaver.EXPECT().
		DeleteCampaign(ctx, campaignId, userId).
		Return(nil)

	err := service.DeleteCampaign(ctx, campaignId, userId)

	require.NoError(t, err)
}

func TestDeleteCampaign_NotMaster(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	err := service.DeleteCampaign(ctx, campaignId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNotMaster, errors.Unwrap(err))
}

func TestDeleteCampaign_CampaignNotFound(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaignPlayers(ctx, campaignId).
		Return([]int64{}, nil)

	mockGameSaver.EXPECT().
		DeleteCampaign(ctx, campaignId, userId).
		Return(models.ErrCampaignNotFound)

	err := service.DeleteCampaign(ctx, campaignId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}
