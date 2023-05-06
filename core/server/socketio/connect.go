package socketio

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func (s *Server) connect(client socketio.Conn) error {
	playerId, err := s.authService.VerifyAuthorization(client.RemoteHeader().Get("Authorization"))
	if err != nil {
		client.Emit(errorEvent, err.Error())
		client.Close()
	}

	client.SetContext(&clientContext{
		playerId,
	})
	log.Println("connected:", client.ID())
	return nil
}
