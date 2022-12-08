package handler

import (
	"net/http"
	"uwwolf/game/core"
	"uwwolf/game/enum"
	"uwwolf/game/types"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func APIRouter() *gin.Engine {
	if router != nil {
		return router
	}

	router := gin.Default()

	router.POST("/game", func(ctx *gin.Context) {
		var setting types.GameSetting
		if err := ctx.ShouldBindJSON(&setting); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// fake game id
		gameID := ""
		core.Manager().AddGame(enum.GameID(gameID), &setting)
		ctx.JSON(http.StatusOK, gin.H{"id": gameID})
	})

	return router
}
