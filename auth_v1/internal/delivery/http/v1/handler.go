package v1

import (
	"github.com/alsadx/GM-Tool/internal/auth"
	"github.com/alsadx/GM-Tool/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	playersService service.Players
	mastersService service.Masters
	tokenManager   auth.TokenManager
}

func NewHandler(playersService service.Players, mastersService service.Masters, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		playersService: playersService,
		mastersService: mastersService,
		tokenManager:   tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup){
	v1 := api.Group("/v1")
	{
		h.initPlayersRoutes(v1)
		h.initMastersRoutes(v1)
	}
}