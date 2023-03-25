package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"uwwolf/app/api"
	"uwwolf/app/service"
	"uwwolf/config"
	"uwwolf/game"
	"uwwolf/storage/postgres"
	"uwwolf/storage/redis"

	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := config.Load(".")

	log.Println("Connecting to Redis...")
	rdb := redis.Connect(config.Redis)

	log.Println("Connecting to PostgreSQL...")
	pdb := postgres.Connect(config.Postgres)

	gameManager := game.Manager(config.Game)

	roomService := service.NewRoomService(rdb)
	gameService := service.NewGameService(config.Game, rdb, pdb, gameManager)
	apiServer := api.NewAPIServer(config.App, roomService, gameService)

	log.Println("Starting API server...")
	svr := apiServer.Server()
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	rdb.Close()
	pdb.Close()

	log.Println("Exited")
}
