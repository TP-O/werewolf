package echo

import (
	"encoding/json"
	"fmt"
	"log"
	"uwwolf/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	*websocket.Upgrader

	boadcaster *Broadcaster
}

func NewHandler(config config.App) *Handler {
	return &Handler{
		Upgrader:   &websocket.Upgrader{},
		boadcaster: NewBroadcaster(),
	}
}

var clientCounter = 0

func (h Handler) connect(ctx *gin.Context) {
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
			dispatcher(client, msg)
		}
	}
}

func (h Handler) Use(router *gin.RouterGroup) {
	router.GET("/", h.connect)
}
