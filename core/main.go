package main

import (
	_ "github.com/lib/pq"
)

func main() {
	// runtime.GOMAXPROCS(1)

	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop()

	// config := config.Load(".")

	// log.Println("Connecting to Redis...")
	// rdb := redis.Connect(config.Redis)
	// defer rdb.Close()
	// log.Println("Connected to Redis...")

	// log.Println("Connecting to PostgreSQL...")
	// pdb := db.Connect(config.Postgres)
	// defer pdb.Close()
	// log.Println("Connected to PostgreSQL...")

	// gameManager := game.NewManager(config.Game)

	// authService, err := service.NewAuthService(config.Firebase)
	// if err != nil {
	// 	panic(err)
	// }

	// roomService := service.NewRoomService(rdb)
	// gameService := service.NewGameService(config.Game, rdb, pdb, gameManager)
	// communicationService := service.NewCommunicationService(config.App.SecretKey)
	// server := server.NewServer(config.App, config.Game, rdb, authService, roomService, gameService, communicationService)

	// go func() {
	// 	log.Printf("Server is listening on port %d", config.App.Port)
	// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Panic(err)
	// 	}
	// }()

	// <-ctx.Done()
	// log.Println("Shutting down...")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := server.Shutdown(ctx); err != nil {
	// 	log.Println(err.Error())
	// }

	// log.Println("Exited")
}

// type Server struct {
// 	*http.Server

// 	config      config.App
// 	roomService service.RoomService
// 	gameService service.GameService
// }

// func NewServer(config config.App, gameCfg config.Game, rdb *redis.ClusterClient, authService service.AuthService, roomService service.RoomService, gameService service.GameService, communicationService service.CommunicationService) *Server {
// 	if config.Env == "production" {
// 		gin.SetMode(gin.ReleaseMode)
// 	}
