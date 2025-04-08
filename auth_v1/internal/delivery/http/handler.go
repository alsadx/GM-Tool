package http

import (
	"log/slog"
	"net/http"

	"github.com/alsadx/GM-Tool/internal/auth"
	v1 "github.com/alsadx/GM-Tool/internal/delivery/http/v1"
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

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.playersService, h.mastersService, h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}

func (h *Handler) Init(host, port string) *gin.Engine {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		v1.LoggerMiddleware(slog.Default()),
	)

	// docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, port)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}
