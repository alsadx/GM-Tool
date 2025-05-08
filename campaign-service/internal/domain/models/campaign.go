package models

import "time"

type Campaign struct{
	Id int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PlayersCount int64 `json:"players_count"`
	PlayersId *[]int64 `json:"players_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CampaignForPlayer struct{
	Id int64 `json:"id"`
	Name string `json:"name"`
	MasterId int64 `json:"master_id"`
}