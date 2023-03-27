package echo

import (
	"log"
	"uwwolf/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	*websocket.Upgrader
}

func NewHandler(config config.App) *Handler {
	svr := &Handler{
		Upgrader: &websocket.Upgrader{},
	}

	return svr
}

func (h Handler) connect(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	c, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv:%s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (h Handler) Use(router *gin.RouterGroup) {
	router.GET("/connect", h.connect)
}
