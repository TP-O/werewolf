package api

import (
	"net/http"
	"uwwolf/app/dto"
	"uwwolf/game/types"

	"github.com/gin-gonic/gin"
)

// updateGameConfig replaces old game config to the new one.
func (s ApiServer) updateGameConfig(ctx *gin.Context) {
	playerID := types.PlayerID(ctx.GetString("playerID"))

	var payload dto.UpdateGameConfigDto
	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	room := s.roomService.PlayerWaitingRoom(playerID)
	if room == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "You're not in any room!",
		})
		return
	}

	if playerID != room.OwnerID {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "Only the room owner can start the game!",
		})
		return
	}

	if err := s.gameService.UpdateGameConfig(room.ID, payload); err != nil {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "Unable to update setting!",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ok",
	})
}
