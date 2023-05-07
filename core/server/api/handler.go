package api

import (
	"uwwolf/config"
	"uwwolf/server/service"
	"uwwolf/server/socketio"
	"uwwolf/util/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	config               config.App
	roomService          service.RoomService
	gameService          service.GameService
	communicationService service.CommunicationService
	socketServer         *socketio.SocketServer
	rdb                  *redis.ClusterClient
}

func NewHandler(config config.App, socketio *socketio.SocketServer, rdb *redis.ClusterClient, roomService service.RoomService, gameService service.GameService, communicationService service.CommunicationService) *Handler {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.Setup(v)
	}

	return &Handler{
		config:               config,
		socketServer:         socketio,
		roomService:          roomService,
		gameService:          gameService,
		communicationService: communicationService,
		rdb:                  rdb,
	}
}

func (h Handler) Use(router *gin.RouterGroup) {
	gameSetup := router.Group("/game")
	gameSetup.Use(h.WaitingRoomOwner)

	gameSetup.POST("/config", h.ReplaceGameConfig)
	gameSetup.POST("/start", h.StartGame)
}
