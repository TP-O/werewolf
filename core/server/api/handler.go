package api

import (
	"uwwolf/config"
	"uwwolf/server/service"
	"uwwolf/util/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	config      config.App
	roomService service.RoomService
	gameService service.GameService
}

func NewHandler(config config.App, roomService service.RoomService, gameService service.GameService) *Handler {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.Setup(v)
	}

	return &Handler{
		config:      config,
		roomService: roomService,
		gameService: gameService,
	}
}

func (h Handler) Use(router *gin.RouterGroup) {
	gameSetup := router.Group("/game")
	gameSetup.Use(h.WaitingRoomOwner)

	gameSetup.POST("/config", h.ReplaceGameConfig)
	gameSetup.POST("/start", h.StartGame)
}
