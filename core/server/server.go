package server

import (
	"fmt"
	"net/http"
	"uwwolf/config"
	"uwwolf/server/api"
	"uwwolf/server/service"
	"uwwolf/server/socketio"
	"uwwolf/util/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	*http.Server

	config      config.App
	roomService service.RoomService
	gameService service.GameService
}

func NewServer(config config.App, gameCfg config.Game, rdb *redis.ClusterClient, authService service.AuthService, roomService service.RoomService, gameService service.GameService, communicationService service.CommunicationService) *Server {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.Setup(v)
	}

	router := gin.Default()
	router.Use(gin.Recovery())

	socketIoServer := socketio.NewServer(gameCfg, authService)
	router.GET("/socket.io/*any", gin.WrapH(socketIoServer))
	router.POST("/socket.io/*any", gin.WrapH(socketIoServer))

	apiRouter := router.Group("/api")
	apiHandler := api.NewHandler(config, socketIoServer, rdb, roomService, gameService, communicationService)
	apiHandler.Use(apiRouter)

	svr := &Server{
		config:      config,
		roomService: roomService,
		gameService: gameService,
	}
	svr.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	return svr
}
