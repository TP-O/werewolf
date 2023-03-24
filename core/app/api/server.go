package api

import (
	"fmt"
	"uwwolf/app/service"
	"uwwolf/app/validation"
	"uwwolf/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ApiServer struct {
	roomService service.RoomService
	gameService service.GameService
}

func NewAPIServer(roomService service.RoomService, gameService service.GameService) *ApiServer {
	return &ApiServer{
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

func (s ApiServer) Run() {
	if config.App().Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.ImproveValidator(v)
	}

	route := s.setupRouter()
	if err := route.Run(fmt.Sprintf(":%v", config.App().Port)); err != nil {
		panic(err)
	}
}
