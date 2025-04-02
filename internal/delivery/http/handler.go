package http

import "github.com/alsadx/GM-Tool/internal/service"

type Handler struct {
	playersService service.Players
	mastersService service.Masters
}

