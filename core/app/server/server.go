package server

import (
	"fmt"
	"net/http"
	"uwwolf/app/server/api"
	ws "uwwolf/app/server/echo"
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

	router := gin.Default()
	router.Use(gin.Recovery())

	apiRouter := router.Group("/api")
	apiHandler := api.NewHandler(config, roomService, gameService)
	apiHandler.Use(apiRouter)

	echoRouter := router.Group("/echo")
	echoHandler := ws.NewHandler(config)
	echoHandler.Use(echoRouter)

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
