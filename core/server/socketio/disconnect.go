package socketio

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func (s *SocketServer) disconnect(client socketio.Conn, reason string) {
	log.Println("closed", reason)
}
