package api

import (
	"context"
	"net/http"
	"uwwolf/server/data"
	"uwwolf/server/enum"
	"uwwolf/server/service"

	"github.com/gin-gonic/gin"
)

// StartGame creates a game moderator and then starts the game.
func (h Handler) StartGame(ctx *gin.Context) {
	v, _ := ctx.Get(enum.WaitingRoomCtxKey)
	room, ok := v.(*data.WaitingRoom)
	if room == nil || !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to start game!",
		})
		return
	}

	gameCfg := h.gameService.GameConfig(room.ID)
	if err := h.gameService.CheckBeforeRegistration(*room, gameCfg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	mod, err := h.gameService.RegisterGame(gameCfg, room.PlayerIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	mod.StartGame()
	set := make([]any, len(room.PlayerIDs)*2, len(room.PlayerIDs)*2)
	for _, id := range room.PlayerIDs {
		set = append(set, id, "in_game")
	}
	h.rdb.MSet(context.Background(), set...)
	h.communicationService.BroadcastToRoom(room.ID, service.CommunicationEventMsg{
		Event:   "start",
		Message: mod.GameID(),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ok",
	})
}
