package api

import (
	"net/http"
	"uwwolf/game/types"
	"uwwolf/server/enum"

	"github.com/gin-gonic/gin"
)

// WaitingRoomOwner gets the waiting room owned by the authenticated player.
func (h Handler) WaitingRoomOwner(ctx *gin.Context) {
	playerID := types.PlayerID(ctx.GetString(enum.PlayerIDCtxKey))

	room, ok := h.roomService.PlayerWaitingRoom(playerID)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "You're not in any room!",
		})
		return
	}

	if playerID != room.OwnerID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Only the room owner can start the game!",
		})
		return
	}

	ctx.Set(enum.WaitingRoomCtxKey, &room)
}
