package socketio

import (
	"context"
	"log"
	"strconv"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/internal/infra/db/redis"

	socketio "github.com/googollee/go-socket.io"
)

func (s *Server) Connect(client socketio.Conn) error {
	playerId, err := s.authService.VerifyAuthorization(client.RemoteHeader().Get("Authorization"))
	if err != nil {
		client.Emit(errorEvent, err.Error())
		client.Close()
	}

	pipe := s.rdb.Pipeline()
	status := pipe.Get(
		context.Background(),
		redis.PlayerStatus+string(playerId),
	)
	sid := pipe.Get(
		context.Background(),
		redis.SocketId+string(playerId),
	)
	pipe.Exec(context.Background())

	if status.String() != "in_game" {
		client.Emit(errorEvent, "Not in any game!")
		client.Close()
	}

	if sid.String() != "" {
		client.Emit(errorEvent, "Someone is using this account!")
		client.Close()
	}

	game, _ := s.db.PlayingGame(context.Background(), string(playerId))

	client.SetContext(&clientContext{
		playerId: playerId,
		gameId:   types.GameID(game.ID),
	})
	client.Join(strconv.Itoa(int(game.ID)))
	log.Println("connected:", client.ID())
	return nil
}
