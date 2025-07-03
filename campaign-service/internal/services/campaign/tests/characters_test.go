package tests

import (
	"campaigntool/internal/domain/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCharacter_HappyPath(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		IsPlayer(ctx, campaignId, userId).
		Return(false, nil)

	mockGameSaver.EXPECT().
		AddCharacter(ctx, campaignId, userId, charId).
		Return(nil)

	err := service.AddCharacter(ctx, campaignId, userId, charId)

	require.NoError(t, err)
}

func TestAddCharacter_NotPlayer(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsPlayer(ctx, campaignId, userId).
		Return(false, nil)

	err := service.AddCharacter(ctx, campaignId, userId, charId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNotPlayer, errors.Unwrap(err))
}

func TestAddCharacter_CharacterInCampaign(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsPlayer(ctx, campaignId, userId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		AddCharacter(ctx, campaignId, userId, charId).
		Return(nil)

	err := service.AddCharacter(ctx, campaignId, userId, charId)

	require.NoError(t, err)
}

func TestRemoveCharacter_HappyPath(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsCharacterOwner(ctx, campaignId, userId, charId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemoveCharacter(ctx, campaignId, userId, charId).
		Return(nil)

	err := service.RemoveCharacter(ctx, campaignId, userId, charId)

	require.NoError(t, err)
}

func TestRemoveCharacter_NotCharacterOwner(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsCharacterOwner(ctx, campaignId, userId, charId).
		Return(false, nil)

	err := service.RemoveCharacter(ctx, campaignId, userId, charId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNotCharacterOwner, errors.Unwrap(err))
}

func TestRemoveCharacter_NotFound(t *testing.T) {
	service, mockGameSaver, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsCharacterOwner(ctx, campaignId, userId, charId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemoveCharacter(ctx, campaignId, userId, charId).
		Return(models.ErrCharacterNotFound)

	err := service.RemoveCharacter(ctx, campaignId, userId, charId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCharacterNotFound, errors.Unwrap(err))
}

func TestGetCampaignCharacters_HappyPath(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaignCharacters(ctx, campaignId).
		Return([]*models.CampaignCharacter{}, nil)

	characters, err := service.GetCampaignCharacters(ctx, campaignId, userId)

	require.NoError(t, err)
	assert.Equal(t, []*models.CampaignCharacter{}, characters)
}

func TestGetCampaignCharacters_NotFound(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaignCharacters(ctx, campaignId).
		Return(nil, models.ErrCampaignNotFound)

	_, err := service.GetCampaignCharacters(ctx, campaignId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}

func TestGetCampaignCharacters_NotMaster(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	_, err := service.GetCampaignCharacters(ctx, campaignId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNotMaster, errors.Unwrap(err))
}

func TestGetPlayerCharacters_HappyPath(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		GetPlayerCharacters(ctx, campaignId, userId).
		Return([]int64{}, nil)

	characters, err := service.GetPlayerCharacters(ctx, campaignId, userId, userId)

	require.NoError(t, err)
	assert.Equal(t, []int64{}, characters)
}

func TestGetPlayerCharacters_HappyPath2(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	playerId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetPlayerCharacters(ctx, campaignId, playerId).
		Return([]int64{}, nil)

	characters, err := service.GetPlayerCharacters(ctx, campaignId, playerId, userId)

	require.NoError(t, err)
	assert.Equal(t, []int64{}, characters)
}

func TestGetPlayerCharacters_NotCharacterOwner(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)
	playerId := int64(111)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	_, err := service.GetPlayerCharacters(ctx, campaignId, playerId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrNotCharacterOwner, errors.Unwrap(err))
}

func TestGetPlayerCharacters_NotFound(t *testing.T) {
	service, _, mockGameProvider := setupTest(t)

	ctx := context.Background()
	campaignId := int64(123)
	userId := int64(123)

	mockGameProvider.EXPECT().
		IsMaster(ctx, campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		GetPlayerCharacters(ctx, campaignId, userId).
		Return(nil, models.ErrCampaignNotFound)

	_, err := service.GetPlayerCharacters(ctx, campaignId, userId, userId)

	require.Error(t, err)
	assert.Equal(t, models.ErrCampaignNotFound, errors.Unwrap(err))
}
