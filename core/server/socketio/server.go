package socketio

import (
	"log"
	"net/http"
	"uwwolf/config"
	"uwwolf/db"
	"uwwolf/game"
	"uwwolf/game/types"
	"uwwolf/server/service"

	socketio "github.com/googollee/go-socket.io"
	"github.com/redis/go-redis/v9"
)

type message[T any] struct {
	Event string `json:"event"`
	Data  T      `json:"data"`
}

type clientContext struct {
	playerId types.PlayerID
	gameId   types.GameID
}

type SocketServer struct {
	*socketio.Server
	authService service.AuthService
	gameManger  game.Manager
	db          db.Store
	rdb         *redis.ClusterClient
}

const defaultNamespace = "/"

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func NewServer(gameConfig config.Game, authService service.AuthService) *SocketServer {
	server := &SocketServer{
		Server:      socketio.NewServer(nil),
		authService: authService,
		gameManger:  game.NewManager(gameConfig),
	}

	server.OnConnect(defaultNamespace, server.connect)
	server.OnDisconnect(defaultNamespace, server.disconnect)
	server.OnError(defaultNamespace, server.handleError)
	server.OnEvent(defaultNamespace, syncPositionEvent, server.syncPosition)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	return server
}
