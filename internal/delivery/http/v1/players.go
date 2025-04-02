package v1

import (
	"net/http"

	"github.com/alsadx/GM-Tool/internal/service"
	"github.com/gin-gonic/gin"
)

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) initPlayersRoutes(api *gin.RouterGroup) {
	players := api.Group("/players")
	{
		players.POST("/sign-up", h.playerSignUp)
		players.POST("/sign-in", h.playerSignIn)
		players.POST("/auth/refresh", h.playerRefresh)
	}
}

func (h *Handler) playerSignUp(c *gin.Context) {
	var inp service.SignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.playersService.SignUp(c.Request.Context(), inp)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) playerSignIn(c *gin.Context) {
	var input service.SignInInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.playersService.SignIn(c.Request.Context(), input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) playerRefresh(c *gin.Context) {
	var input refreshInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.playersService.RefreshTokens(c.Request.Context(), input.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}