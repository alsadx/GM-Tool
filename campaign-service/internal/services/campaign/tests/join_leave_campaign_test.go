package tests

import (
	"campaigntool/internal/domain/models"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinCampaign_Success(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	charId := int64(111)
	userId := int64(123)

	mockGameProvider.EXPECT().
		CheckInviteCode(ctx, gomock.Any()).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(ctx, campaignId, userId, charId).
		Return(nil)

	err := service.JoinCampaign(ctx, "ABCDEF", userId, charId)

	require.NoError(t, err)
}

func TestJoinCampaign_CampaignNotFound(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	charId := int64(111)
	userId := int64(123)

	mockGameProvider.EXPECT().
		CheckInviteCode(ctx, gomock.Any()).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(ctx, campaignId, userId, charId).
		Return(models.ErrCampaignNotFound)

	err := service.JoinCampaign(ctx, "ABCDEF", userId, charId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}

func TestJoinCampaign_InvalidCode(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()

	charId := int64(111)
	userId := int64(123)

	mockGameProvider.EXPECT().
		CheckInviteCode(ctx, gomock.Any()).
		Return(int64(0), models.ErrInvalidCode)

	err := service.JoinCampaign(ctx, "ABCDEF", userId, charId)

	require.Error(t, err)
	assert.Equal(t, models.ErrInvalidCode, errors.Unwrap(err))
}

func TestLeaveCampaign_HappyPath(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsPlayer(ctx, campaignId, userId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemovePlayer(ctx, campaignId, userId).
		Return(nil)

	err := service.LeaveCampaign(ctx, campaignId, userId)

	require.NoError(t, err)
}

func TestLeaveCampaign_CampaignNotFound(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsPlayer(ctx, campaignId, userId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemovePlayer(ctx, campaignId, userId).
		Return(models.ErrCampaignNotFound)

	err := service.LeaveCampaign(ctx, campaignId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}

func TestLeaveCampaign_NotPlayer(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsPlayer(ctx, campaignId, userId).
		Return(false, nil)

	err := service.LeaveCampaign(ctx, campaignId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNotPlayer, errors.Unwrap(err))
}