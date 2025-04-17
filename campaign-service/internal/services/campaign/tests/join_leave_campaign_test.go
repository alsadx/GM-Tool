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

	ctx := context.WithValue(context.Background(), "user_id", 1)
	campaignId := int32(123)

	mockGameProvider.EXPECT().
		CheckInviteCode(ctx, gomock.Any()).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(ctx, campaignId, 1).
		Return(nil)

	err := service.JoinCampaign(ctx, "ABCDEF", 1)

	require.NoError(t, err)
}

func TestJoinCampaign_CampaignNotFound(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)
	campaignId := int32(123)

	mockGameProvider.EXPECT().
		CheckInviteCode(ctx, gomock.Any()).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(ctx, campaignId, 1).
		Return(models.ErrCampaignNotFound)

	err := service.JoinCampaign(ctx, "ABCDEF", 1)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}

func TestJoinCampaign_InvalidCode(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)

	mockGameProvider.EXPECT().
		CheckInviteCode(ctx, gomock.Any()).
		Return(int32(0), models.ErrInvalidCode)

	err := service.JoinCampaign(ctx, "ABCDEF", 1)

	require.Error(t, err)
	assert.Equal(t, models.ErrInvalidCode, errors.Unwrap(err))
}

func TestLeaveCampaign_HappyPath(t *testing.T) {
	service, mockGameSaver, _ := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)
	campaignId := int32(123)

	mockGameSaver.EXPECT().
		RemovePlayer(ctx, campaignId, 1).
		Return(nil)

	err := service.LeaveCampaign(ctx, campaignId, 1)

	require.NoError(t, err)
}

func TestLeaveCampaign_CampaignNotFound(t *testing.T) {
	service, mockGameSaver, _ := setupTest(t)

	ctx := context.WithValue(context.Background(), "user_id", 1)
	campaignId := int32(123)

	mockGameSaver.EXPECT().
		RemovePlayer(ctx, campaignId, 1).
		Return(models.ErrCampaignNotFound)

	err := service.LeaveCampaign(ctx, campaignId, 1)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}
