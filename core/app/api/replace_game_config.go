package api

import (
	"fmt"
	"net/http"
	"uwwolf/app/data"
	"uwwolf/app/dto"
	"uwwolf/app/enum"
	"uwwolf/app/validation"

	"github.com/gin-gonic/gin"
)

// ReplaceGameConfig replaces old game config to the new one.
func (s Server) ReplaceGameConfig(ctx *gin.Context) {
	var payload dto.ReplaceGameConfigDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println(validation.FormatValidationError(err))

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request!",
			"errors":  validation.FormatValidationError(err),
		})
		return
	}

	v, _ := ctx.Get(enum.WaitingRoomCtxKey)
	room, ok := v.(*data.WaitingRoom)
	if room == nil || !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to update game config!",
		})
		return
	}

	if err := s.gameService.UpdateGameConfig(room.ID, payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ok",
	})
}
