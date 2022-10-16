package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"

	"uwwolf/app/enum"
	"uwwolf/app/instance"
	"uwwolf/app/service"
	"uwwolf/config"
)

func StartSocketIO() {
	defer instance.SocketIOServer.Close()

	instance.SocketIOServer.OnConnect("/", func(client socketio.Conn) error {
		if playerId, err := service.Verify(client.RemoteHeader().Get("Authorization")); err != nil {
			return err
		} else {
			gameId, err := instance.RedisClient.Get(
				context.Background(),
				enum.PlayerId2GameIdCacheNamespace+playerId,
			).Result()

			if err != nil {
				return errors.New("Start game before connecting to server!")
			}

			_, err = instance.RedisClient.Get(
				context.Background(),
				enum.PlayerId2GameIdCacheNamespace+playerId,
			).Result()

			if err == nil {
				return errors.New("Someone is playing!")
			}

			redisPipe := instance.RedisClient.Pipeline()
			redisPipe.Set(
				context.Background(),
				enum.PlayerId2SocketIdCacheNamespace+playerId,
				client.ID(),
				-1,
			)
			redisPipe.MSet(
				context.Background(),
				enum.SocketIdCacheNamespace+client.ID(),
				"playerId",
				playerId,
				"gameId",
				gameId,
				-1,
			)
			redisPipe.Exec(context.Background())
			client.Join(gameId)
		}

		return nil
	})

	instance.SocketIOServer.OnDisconnect("/", func(client socketio.Conn, reason string) {
		a := instance.RedisClient.HGetAll(
			context.Background(),
			enum.SocketIdCacheNamespace+client.ID(),
		).Val()
		playerId := a["playerId"]
		gameId := a["gameId"]
		redisPipe := instance.RedisClient.Pipeline()
		redisPipe.Del(context.Background(), enum.PlayerId2SocketIdCacheNamespace+playerId)
		redisPipe.Del(context.Background(), enum.SocketIdCacheNamespace+client.ID())
		redisPipe.Exec(context.Background())

		instance.SocketIOServer.BroadcastToRoom("", gameId, "leave", fiber.Map{
			"id":      playerId,
			"message": "Leave room",
		})
	})

	instance.SocketIOServer.OnError("/", func(client socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	instance.SocketIOServer.OnEvent("/", "notice", func(client socketio.Conn, msg string) {
		log.Println("notice:", msg)
		client.Emit("reply", "have "+msg)
	})

	go func() {
		if err := instance.SocketIOServer.Serve(); err != nil {
			log.Fatalf("Socketio listen error: %s\n", err)
		}
	}()

	http.Handle("/socket.io/", instance.SocketIOServer)

	log.Println("SocketIO Server is running at http://127.0.0.1:" + strconv.Itoa(config.App.WsPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.App.WsPort), nil))
}
