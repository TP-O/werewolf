package echo

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Room struct {
	clientIDs []string
}

type Broadcaster struct {
	clients map[string]*websocket.Conn
	rooms   map[string]*Room
	events  map[string]func(client *websocket.Conn)
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{}
}

func (b Broadcaster) AddClient(clientID string, client *websocket.Conn) {
	b.clients[clientID] = client

	for {
		var msg any
		err := client.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv:%s", msg)
	}
}

func (b Broadcaster) Emit(msg string, to []string) {
	for _, clientID := range to {
		err := b.clients[clientID].WriteMessage(1, []byte(msg))
		if err != nil {
			fmt.Println("write:", err)
		}
	}
}

func (b Broadcaster) EmitRoom(msg string, roomID string) error {
	room := b.rooms[roomID]
	if room == nil {
		return fmt.Errorf("Room does not exist!")
	}

	b.Emit(msg, room.clientIDs)
	return nil
}

func (b Broadcaster) On(event string, fn func(client *websocket.Conn)) {
	b.events[event] = fn
}
