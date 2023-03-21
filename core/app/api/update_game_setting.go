package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"uwwolf/app/dto"
	"uwwolf/app/service"
	"uwwolf/db/rdb"
	"uwwolf/game/types"

	"github.com/gin-gonic/gin"
)

func (s ApiServer) updateGameSetting(ctx *gin.Context) {
	playerID := types.PlayerID(ctx.GetString("playerID"))

	var payload dto.UpdateGameSettingDto
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

	room.DiscussionDuration = payload.DiscussionDuration
	room.TurnDuration = payload.TurnDuration
	room.RoleIDs = payload.RoleIDs
	room.RequiredRoleIDs = payload.RequiredRoleIDs
	room.NumberWerewolves = payload.NumberWerewolves

	roomByte, err := json.Marshal(room)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something wrong",
		})
		return
	}

	rdb.Client().Set(context.Background(), fmt.Sprintf("%v:%v", service.WaitingRoomRedisNamespace, room.ID), string(roomByte), -1)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ok",
	})
}
