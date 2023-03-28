package echo

type Event = int

const (
	SyncEvent Event = iota
)

type EventMessage struct {
	Event Event `json:"event"`
}

type Dispatcher = func(client *Client, msg EventMessage)

type Broadcaster struct {
	clients     map[ClientID]*Client
	rooms       map[RoomID]*Room
	dispatchers map[Event]Dispatcher
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients:     make(map[ClientID]*Client),
		rooms:       make(map[RoomID]*Room),
		dispatchers: make(map[Event]Dispatcher),
	}
}

func (b Broadcaster) Client(clientID ClientID) *Client {
	return b.clients[clientID]
}

func (b *Broadcaster) AddClient(client *Client) {
	if client != nil {
		b.clients[client.id] = client
	}
}

func (b *Broadcaster) RemoveClient(clientID ClientID) {
	client := b.clients[clientID]
	if client == nil {
		return
	}

	client.LeaveAllRooms()
	delete(b.clients, clientID)
}

func (b Broadcaster) Emit(msg EventMessage, to []ClientID) {
	for _, clientID := range to {
		b.clients[clientID].conn.WriteJSON(msg) // nolint errcheck
	}
}

func (b Broadcaster) EmitRoom(msg EventMessage, to []RoomID) {
	for _, roomID := range to {
		room := b.rooms[roomID]
		if room != nil {
			room.Broadcast(msg)
		}
	}
}

func (b *Broadcaster) On(event Event, dispatcher Dispatcher) {
	b.dispatchers[event] = dispatcher
}

func (b *Broadcaster) Dispatcher(event Event) Dispatcher {
	return b.dispatchers[event]
}
