package models

type Player struct {
    ID   int64  `json:"user_id"`
    Name string `json:"name"`
}

type CampaignResponse struct {
    ID          int64    `json:"campaign_id"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    CreatedAt   string   `json:"created_at"`
    Players     []Player `json:"players"`
}