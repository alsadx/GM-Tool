package v1

import (
	"net/http"

	"github.com/alsadx/GM-Tool/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initMastersRoutes(api *gin.RouterGroup) {
	masters := api.Group("/masters")
	{
		masters.POST("/sign-up", h.masterSignUp)
		masters.POST("/sign-in", h.masterSignIn)
		masters.POST("/auth/refresh", h.masterRefresh)
	}
}

func (h *Handler) masterSignUp(c *gin.Context) {
	var inp service.SignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.mastersService.SignUp(c.Request.Context(), inp)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) masterSignIn(c *gin.Context) {
	var input service.SignInInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.mastersService.SignIn(c.Request.Context(), input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) masterRefresh(c *gin.Context) {
	var input refreshInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.mastersService.RefreshTokens(c.Request.Context(), input.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}