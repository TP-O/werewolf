package socketio

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func (s *Server) connect(client socketio.Conn) error {
	fmt.Println("CCC")
	client.SetContext(&clientContext{
		playerID: "1",
	})
	log.Println("connected:", client.ID())
	return nil
}
