package socketio

import (
	"log"
	game "uwwolf/internal/app/game/logic"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/internal/config"
	db "uwwolf/internal/infra/db/postgres"
	"uwwolf/internal/infra/firebase"
	"uwwolf/internal/port/driver"

	socketio "github.com/googollee/go-socket.io"
	"github.com/redis/go-redis/v9"
)

type message[T any] struct {
	Event string `json:"event"`
	Data  T      `json:"data"`
}

type clientContext struct {
	playerId types.PlayerId
	gameId   types.GameID
}

type Server struct {
	*socketio.Server
	authService firebase.AuthService
	gameManger  contract.Manager
	db          db.Store
	rdb         *redis.ClusterClient
}

const defaultNamespace = "/"

// var allowOriginFunc = func(r *http.Request) bool {
// 	return true
// }

func NewServer(gameConfig config.Game, authService firebase.AuthService) driver.SocketIoPort {
	server := &Server{
		Server:      socketio.NewServer(nil),
		authService: authService,
		gameManger:  game.NewManager(gameConfig),
	}

	server.OnConnect(defaultNamespace, server.Connect)
	server.OnDisconnect(defaultNamespace, server.Disconnect)
	server.OnError(defaultNamespace, server.HandleError)
	server.OnEvent(defaultNamespace, syncPositionEvent, server.SyncPosition)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	return server
}
