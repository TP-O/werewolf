package api

// WaitingRoomOwner gets the waiting room owned by the authenticated player.
// func (s Server) WaitingRoomOwner(ctx *gin.Context) {
// 	playerID := types.PlayerId(ctx.GetString(enum.PlayerIDCtxKey))

// 	room, ok := s.roomService.PlayerWaitingRoom(playerID)
// 	if !ok {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"message": "You're not in any room!",
// 		})
// 		return
// 	}

// 	if playerID != room.OwnerID {
// 		ctx.JSON(http.StatusForbidden, gin.H{
// 			"message": "Only the room owner can start the game!",
// 		})
// 		return
// 	}

// 	ctx.Set(enum.WaitingRoomCtxKey, &room)
// }
