package api

import (
	"uwwolf/internal/config"
	"uwwolf/internal/port/driver"
	"uwwolf/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config       config.App
	socketIoPort driver.SocketIoPort
}

func NewHandler(config config.App, socketIo driver.SocketIoPort) driver.HttpApiPort {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.Setup(v)
	}

	return &Server{
		config:       config,
		socketIoPort: socketIo,
	}
}

func (h Server) Use(router *gin.RouterGroup) {
	gameSetup := router.Group("/game")
	// gameSetup.Use(h.WaitingRoomOwner)

	gameSetup.POST("/setting", h.UpdateGameSetting)
	gameSetup.POST("/start", h.StartGame)
}
