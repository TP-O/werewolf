package api

import (
	"github.com/gin-gonic/gin"
)

// ReplaceGameConfig replaces old game config to the new one.
func (s Server) UpdateGameSetting(ctx *gin.Context) {
	// var payload dto.ReplaceGameConfigDto
	// if err := ctx.ShouldBindJSON(&payload); err != nil {
	// 	fmt.Println(validation.FormatValidationError(err))

	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Invalid request!",
	// 		"errors":  validation.FormatValidationError(err),
	// 	})
	// 	return
	// }

	// v, _ := ctx.Get(enum.WaitingRoomCtxKey)
	// room, ok := v.(*data.WaitingRoom)
	// if room == nil || !ok {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Unable to update game config!",
	// 	})
	// 	return
	// }

	// if err := h.gameService.UpdateGameConfig(room.ID, payload); err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Something went wrong!",
	// 	})
	// 	return
	// }

	// h.communicationService.BroadcastToRoom(room.ID, service.CommunicationEventMsg{
	// 	Event:   "room_setting",
	// 	Message: payload,
	// })

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "Ok",
	// })
}
