package server

import (
	"log"
	"net/http"
	"strconv"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"uwwolf/config"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func StartSocketIO() {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	defer server.Close()

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socketio listen error: %s\n", err)
		}
	}()

	http.Handle("/socket.io/", server)

	log.Println("SocketIO Server is running at http://127.0.0.1:" + strconv.Itoa(config.App.WsPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.App.WsPort), nil))
}
