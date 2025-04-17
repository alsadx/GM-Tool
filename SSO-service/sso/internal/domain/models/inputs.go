package models

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"` // TODO: add min length
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"` // TODO: add min length
	Name     string `json:"name" validate:"required,min=3"`
}

type IsAdminInput struct {
	UserId int64 `json:"user_id" validate:"required,min=1"`
}