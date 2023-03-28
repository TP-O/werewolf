package echo

import "golang.org/x/exp/maps"

type RoomID = string

type Room struct {
	id          RoomID
	broadcaster *Broadcaster
	clients     map[ClientID]*Client
}

func NewRoom(id RoomID, broadcaster *Broadcaster) *Room {
	return &Room{
		id:          id,
		broadcaster: broadcaster,
		clients:     make(map[string]*Client),
	}
}

func (r *Room) AddClient(client *Client) {
	r.broadcaster.AddClient(client)
	r.clients[client.id] = client
}

func (r *Room) RemoveClient(clientID ClientID) {
	delete(r.clients, clientID)

	if len(r.clients) == 0 {
		delete(r.broadcaster.rooms, r.id)
	}
}

func (r Room) Broadcast(msg EventMessage) {
	r.broadcaster.Emit(msg, maps.Keys(r.clients))
}
