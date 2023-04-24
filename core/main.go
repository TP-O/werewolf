package main

import (
	"fmt"
	"uwwolf/pkg/quadtree"

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
	// pdb := postgres.Connect(config.Postgres)
	// defer pdb.Close()
	// log.Println("Connected to PostgreSQL...")

	// gameManager := game.Manager(config.Game)

	// roomService := service.NewRoomService(rdb)
	// gameService := service.NewGameService(config.Game, rdb, pdb, gameManager)
	// server := server.NewServer(config.App, config.Game, roomService, gameService)

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

	// m := tool.NewMap()
	// m.AddEntity(1, tool.EntitySettings{
	// 	X:     2,
	// 	Y:     2,
	// 	Speed: 3,
	// })
	// if ok, err := m.MoveEntity(1, orb.Point{2.1, 2.1}); !ok {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Moved")
	// }

	qt := &quadtree.Quadtree{
		Bounds: quadtree.Bounds{
			X:      0,
			Y:      0,
			Width:  100,
			Height: 100,
		},
		MaxObjects: 10,
		MaxLevels:  4,
		Level:      0,
		Objects:    make([]quadtree.IBounds, 0),
		Nodes:      make([]quadtree.Quadtree, 0),
	}

	qt.Insert(&quadtree.Bounds{
		X:      1,
		Y:      1,
		Width:  10,
		Height: 10,
	})
	qt.Insert(&quadtree.Bounds{
		X:      1,
		Y:      1,
		Width:  10,
		Height: 10,
	})
	qt.Insert(&quadtree.Bounds{
		X:      1,
		Y:      1,
		Width:  10,
		Height: 10,
	})
	qt.Insert(&quadtree.Bounds{
		X:      50,
		Y:      50,
		Width:  10,
		Height: 10,
	})

	intersections := qt.RetrieveIntersections(&quadtree.Bounds{
		X:      0,
		Y:      0,
		Width:  20,
		Height: 20,
	})

	fmt.Println(intersections)
	fmt.Println(qt.TotalNodes())
}
