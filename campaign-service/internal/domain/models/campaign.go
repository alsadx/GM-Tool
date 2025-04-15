package models

import "time"

type Campaign struct{
	Id int32 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PlayerCount int32 `json:"player_count"`
	CreatedAt time.Time `json:"created_at"`
}

type CampaignForPlayer struct{
	Id int32 `json:"id"`
	Name string `json:"name"`
}