package api

import (
	"fmt"
	"net/http"
	"uwwolf/app/service"
	"uwwolf/app/validation"
	"uwwolf/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	*http.Server

	config      config.App
	roomService service.RoomService
	gameService service.GameService
}

func NewServer(config config.App, roomService service.RoomService, gameService service.GameService) *Server {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.Setup(v)
	}

	svr := &Server{
		config:      config,
		roomService: roomService,
		gameService: gameService,
	}
	svr.Server = &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Port),
		Handler: svr.router(),
	}

	return svr
}

func (s Server) router() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Recovery())

	gameSetup := r.Group("/game")
	gameSetup.Use(s.WaitingRoomOwner)

	gameSetup.POST("/config", s.ReplaceGameConfig)
	gameSetup.POST("/start", s.StartGame)

	return r
}
