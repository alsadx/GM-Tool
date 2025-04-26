package models

import "time"

type Session struct {
	RefreshToken string
	ExpiresAt    time.Time
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
