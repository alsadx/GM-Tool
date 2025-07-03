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

func TestGenerateInviteCode_Success(t *testing.T) {
	service, mockGameSaver, _ := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameSaver.EXPECT().
		SetInviteCode(ctx, campaignId, userId, gomock.Any()).
		Return(nil)

	inviteCode, err := service.GenerateInviteCode(ctx, campaignId, userId)

	require.NoError(t, err)
	assert.NotEmpty(t, inviteCode)
	require.Condition(t, func() bool { return len(inviteCode) == 6 })
}

func TestGenerateInviteCode_CampaignNotFound(t *testing.T) {
	service, mockGameSaver, _ := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameSaver.EXPECT().
		SetInviteCode(ctx, campaignId, userId, gomock.Any()).
		Return(models.ErrCampaignNotFound)

	inviteCode, err := service.GenerateInviteCode(ctx, campaignId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
	assert.Empty(t, inviteCode)
}
