package api

import (
	"fmt"
	"uwwolf/app/service"
	"uwwolf/util"

	"github.com/gin-gonic/gin"
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

	gameGroup := r.Group("/game")
	gameGroup.PUT("/setting", s.registerGame)
	gameGroup.POST("/start", s.updateGameSetting)

	return r
}

func (s ApiServer) Run() {
	r := s.setupRouter()

	if err := r.Run(fmt.Sprintf(":%v", util.Config().App.Port)); err != nil {
		panic(err)
	}
}
