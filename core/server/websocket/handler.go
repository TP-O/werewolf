package echo

import (
	"encoding/json"
	"fmt"
	"log"
	"uwwolf/config"
	"uwwolf/game"
	"uwwolf/game/types"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	*websocket.Upgrader

	boadcaster *Broadcaster
	gameManger game.Manager
}

func NewHandler(config config.App, gConfig config.Game) *Handler {
	h := &Handler{
		Upgrader:   &websocket.Upgrader{},
		gameManger: game.NewManager(gConfig),
	}
	boadcaster := NewBroadcaster()
	boadcaster.On(SyncEvent, h.SyncPosition)
	h.boadcaster = boadcaster

	return h
}

var clientCounter = 0
var gameID = 1

func (h *Handler) connect(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Uprage to WS:", err)
		return
	}
	defer conn.Close()

	client := NewClient(fmt.Sprintf("%d", clientCounter), h.boadcaster, conn)
	h.boadcaster.AddClient(client)
	defer h.boadcaster.RemoveClient(client.id)
	client.JoinRoom(fmt.Sprintf("%d", gameID))
	clientCounter++
	for {
		_, byteMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read message:", err)
			break
		}

		var msg EventMessage
		err = json.Unmarshal(byteMsg, &msg)
		if err != nil {
			continue
		}

		if dispatcher := h.boadcaster.Dispatcher(msg.Event); dispatcher != nil {
			dispatcher(client, msg.Data)
		}
	}
}

func (h *Handler) Use(router *gin.RouterGroup) {
	router.GET("/", h.connect)
}

func (h *Handler) SyncPosition(client *Client, data map[string]any) {
	x, xOk := data["x"].(float64)
	y, yOk := data["y"].(float64)
	if !xOk || !yOk {
		return
	}

	mod := h.gameManger.ModeratorOfPlayer(types.PlayerID(client.ID()))
	if mod != nil {
		ok, _ := mod.MovePlayer(types.PlayerID(client.ID()), x, y)
		if ok {
			client.EmitRoom(EventMessage{
				Event: SyncEvent,
				Data: map[string]any{
					"x":        x,
					"y":        y,
					"playerID": client.ID(),
				},
			}, []RoomID{fmt.Sprintf("%v", mod.GameID())})
		}
	}
}
