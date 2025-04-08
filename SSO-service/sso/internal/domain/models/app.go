package models

type App struct {
	Id   int  `json:"id"`
	Name string `json:"name"`
	SigningKey string `json:"signingKey"`
}