package tests

import (
	"campaigntool/internal/domain/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCampaign_Success(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := "New Test Campaign"
	newDescription := "This is a new test campaign"

	updatedCampaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        newName,
		Description: newDescription,
		MasterId:    userId,
	}

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(ctx, campaignId).
		Return(campaign, nil)

	mockGameSaver.EXPECT().
		UpdateCampaign(ctx, updatedCampaign).
		Return(nil)

	err := service.UpdateCampaign(ctx, campaignId, newName, newDescription, userId)

	require.NoError(t, err)
}

func TestUpdateCampaign_NoChanges(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	
	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        "Test Campaign",
		Description: "This is a test campaign",
		MasterId:    userId,
	}

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(ctx, campaignId).
		Return(campaign, nil)

	err := service.UpdateCampaign(ctx, campaignId, "", "", userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNoChanges, errors.Unwrap(err))
}

func TestUpdateCampaign_NotMaster(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	name := "New Test Campaign"
	description := "This is a new test campaign"
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	err := service.UpdateCampaign(ctx, campaignId, name, description, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNotMaster, errors.Unwrap(err))
}

func TestUpdateCampaign_CampaignNotFound(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	name := "New Test Campaign"
	description := "This is a new test campaign"
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(ctx, campaignId).
		Return(nil, models.ErrCampaignNotFound)

	err := service.UpdateCampaign(ctx, campaignId, name, description, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}
