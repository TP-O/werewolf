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

type ApiServer struct {
	config      config.App
	roomService service.RoomService
	gameService service.GameService
}

func NewAPIServer(config config.App, roomService service.RoomService, gameService service.GameService) *ApiServer {
	return &ApiServer{
		config,
		roomService,
		gameService,
	}
}

func (s *ApiServer) setupRouter() *gin.Engine {
	r := gin.Default()

	gameSetup := r.Group("/game")

	gameSetup.POST("/config", s.ReplaceGameConfig)
	gameSetup.POST("/start", s.StartGame)

	return r
}

func (as ApiServer) Server() *http.Server {
	if as.config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.Setup(v)
	}

	route := as.setupRouter()

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", as.config.Port),
		Handler: route,
	}
}
