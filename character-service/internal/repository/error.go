package repository

import "errors"

var (
	ErrCharacterNotFound = errors.New("character not found")
	ErrEmptyID           = errors.New("empty character ID")
)
