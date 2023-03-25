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
	"uwwolf/db/postgres"
	"uwwolf/db/redis"
	"uwwolf/game"

	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := config.Load(".")

	log.Println("Connecting to Redis...")
	rdb := redis.Connect(config.Redis)
	defer rdb.Close()
	log.Println("Connected to Redis...")

	log.Println("Connecting to PostgreSQL...")
	pdb := postgres.Connect(config.Postgres)
	defer pdb.Close()
	log.Println("Connected to PostgreSQL...")

	gameManager := game.Manager(config.Game)

	roomService := service.NewRoomService(rdb)
	gameService := service.NewGameService(config.Game, rdb, pdb, gameManager)
	apiServer := api.NewAPIServer(config.App, roomService, gameService)
	svr := apiServer.Server()

	go func() {
		log.Println("Starting API server...")
		if err := svr.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Println(err)
	}

	log.Println("Exited")
}
