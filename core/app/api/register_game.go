package api

import (
	"net/http"
	"uwwolf/game/types"

	"github.com/gin-gonic/gin"
)

// registerGame creates a game moderator and then starts the game.
func (s ApiServer) registerGame(ctx *gin.Context) {
	playerID := types.PlayerID(ctx.GetString("playerID"))

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

	mod, err := s.gameService.RegisterGame(&types.ModeratorInit{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	mod.StartGame()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ok",
	})
}
