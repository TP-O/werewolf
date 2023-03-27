package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"uwwolf/app/server"
	"uwwolf/app/service"
	"uwwolf/config"
	"uwwolf/db/postgres"
	"uwwolf/db/redis"
	"uwwolf/game"

	_ "github.com/lib/pq"
)

func main() {
	// runtime.GOMAXPROCS(1)

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
	server := server.NewServer(config.App, roomService, gameService)

	go func() {
		log.Printf("Server is listening on port %d", config.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println(err.Error())
	}

	log.Println("Exited")
}
