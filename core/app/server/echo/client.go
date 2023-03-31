package echo

import (
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

type ClientID = string

type Client struct {
	id          ClientID
	broadcaster *Broadcaster
	conn        *websocket.Conn
	rooms       map[RoomID]*Room
}

func NewClient(id ClientID, boadcaster *Broadcaster, conn *websocket.Conn) *Client {
	return &Client{
		id:          id,
		broadcaster: boadcaster,
		conn:        conn,
		rooms:       make(map[RoomID]*Room),
	}
}

func (c Client) ID() ClientID {
	return c.id
}

func (c Client) EmitRoom(msg EventMessage, to []RoomID) {
	for _, roomID := range to {
		room := c.rooms[roomID]
		if room != nil {
			c.broadcaster.Emit(
				msg,
				lo.Filter(
					maps.Keys(room.clients), func(id string, _ int) bool {
						return id != c.id
					},
				),
			)
		}
	}
}

func (c *Client) JoinRoom(roomID RoomID) {
	room := c.broadcaster.rooms[roomID]
	if room == nil {
		room = NewRoom(roomID, c.broadcaster)
		c.broadcaster.rooms[roomID] = room
	}

	c.rooms[roomID] = room
	room.AddClient(c)
}

func (c *Client) LeaveRoom(roomID RoomID) {
	room := c.rooms[roomID]
	if room != nil {
		room.RemoveClient(c.id)
		delete(c.rooms, roomID)
	}
}

func (c *Client) LeaveAllRooms() {
	for roomID, room := range c.rooms {
		room.RemoveClient(c.id)
		delete(c.rooms, roomID)
	}
}
