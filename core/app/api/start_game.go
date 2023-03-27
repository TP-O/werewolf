package api

import (
	"net/http"
	"uwwolf/app/data"
	"uwwolf/app/enum"

	"github.com/gin-gonic/gin"
)

// StartGame creates a game moderator and then starts the game.
func (s Server) StartGame(ctx *gin.Context) {
	v, _ := ctx.Get(enum.WaitingRoomCtxKey)
	room, ok := v.(*data.WaitingRoom)
	if room == nil || !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to start game!",
		})
		return
	}

	gameCfg := s.gameService.GameConfig(room.ID)
	if err := s.gameService.CheckBeforeRegistration(*room, gameCfg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	mod, err := s.gameService.RegisterGame(gameCfg, room.PlayerIDs)
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
