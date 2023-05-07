package socketio

import (
	"context"
	"log"
	"uwwolf/server/enum"

	socketio "github.com/googollee/go-socket.io"
)

func (s *Server) connect(client socketio.Conn) error {
	playerId, err := s.authService.VerifyAuthorization(client.RemoteHeader().Get("Authorization"))
	if err != nil {
		client.Emit(errorEvent, err.Error())
		client.Close()
	}

	pipe := s.rdb.Pipeline()
	status := pipe.Get(
		context.Background(),
		enum.PlayerStatusNs+string(playerId),
	)
	sid := pipe.Get(
		context.Background(),
		enum.SocketIdNs+string(playerId),
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

	client.SetContext(&clientContext{
		playerId,
	})
	log.Println("connected:", client.ID())
	return nil
}
