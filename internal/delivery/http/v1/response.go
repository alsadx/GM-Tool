package v1

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, status int, message string) {
	logger, ok := c.Request.Context().Value("logger").(*slog.Logger)
	if !ok || logger == nil {
		panic("logger not found in context")
	}

	logger.Error(message)
	c.JSON(status, Response{
		Message: message,
	})
}
