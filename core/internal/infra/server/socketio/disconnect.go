package socketio

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func (s *Server) Disconnect(client socketio.Conn, reason string) {
	log.Println("closed", reason)
}
