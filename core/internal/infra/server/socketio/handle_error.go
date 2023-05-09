package socketio

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

const errorEvent = "error"

func (s *Server) HandleError(client socketio.Conn, e error) {
	client.Emit(errorEvent, e.Error())
	log.Println("meet error:", e)
}
