package models

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash string `json:"password"`
	Name     string `json:"name"`
}