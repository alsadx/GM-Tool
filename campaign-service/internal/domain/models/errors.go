package models

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCampaignNotFound   = errors.New("campaign not found")
	ErrCampaignExists     = errors.New("campaign with this name already exists")
	ErrInvalidCode        = errors.New("invalid invite code")
	ErrPlayerInCampaign   = errors.New("player already in campaign")
	ErrNoCampaigns        = errors.New("campaigns not found")
	ErrCharacterNotFound  = errors.New("character not found")
	ErrCharacterInCampaign = errors.New("character already in campaign")
	ErrNoChanges          = errors.New("no changes detected")
	ErrNotMaster          = errors.New("user is not master")
	ErrNotPlayer          = errors.New("user is not player")
	ErrNotCharacterOwner  = errors.New("user is not character owner")
)
