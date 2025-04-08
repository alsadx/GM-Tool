package v1

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "logger", logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
