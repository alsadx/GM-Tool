package models

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCampaignNotFound   = errors.New("campaign not found")
	ErrCampaignExists     = errors.New("campaign with this name already exists")
	ErrInvalidCode        = errors.New("invalid invite code")
	ErrPlayerInCampaign   = errors.New("player already in campaign")
	ErrNoCampaigns        = errors.New("campaigns not found")
)
