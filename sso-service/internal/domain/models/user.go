package models

// type User struct {
// 	Id       int64  `json:"id"`
// 	Email    string `json:"email"`
// 	PassHash []byte `json:"password"`
// 	Name     string `json:"name"`
// }

type User struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	PassHash []byte `json:"password"`
	Name  string `json:"name"`
	FullName string `json:"full_name"`
	IsAdmin bool `json:"admin"`
	AvatarUrl string `json:"avatar"`
}