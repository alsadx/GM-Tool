package tests

import (
	"campaigntool/internal/domain/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemovePlayer_HappyPath(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	playerId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemovePlayer(ctx, campaignId, playerId).
		Return(nil)

	err := service.RemovePlayer(ctx, campaignId, playerId, userId)

	require.NoError(t, err)
}

func TestRemovePlayer_NotMaster(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	playerId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	err := service.RemovePlayer(ctx, campaignId, playerId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNotMaster, errors.Unwrap(err))
}

func TestRemovePlayer_CampaignNotFound(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	playerId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemovePlayer(ctx, campaignId, playerId).
		Return(models.ErrCampaignNotFound)

	err := service.RemovePlayer(ctx, campaignId, playerId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}

func TestGetCampaignPlayers_HappyPath(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)

	mockGameProvider.EXPECT().
		GetCampaignPlayers(ctx, campaignId).
		Return([]int64{1, 2, 3}, nil)

	players, err := service.GetCampaignPlayers(ctx, campaignId)

	require.NoError(t, err)
	assert.Equal(t, 3, len(players))
}

func TestGetCampaignPlayers_NotFound(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)

	mockGameProvider.EXPECT().
		GetCampaignPlayers(ctx, campaignId).
		Return(nil, models.ErrCampaignNotFound)

	players, err := service.GetCampaignPlayers(ctx, campaignId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
	assert.Nil(t, players)
}