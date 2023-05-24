package game

import (
	roomModel "uwwolf/internal/domain/room/model"
)

// StartGame creates a game moderator and then starts the game.
func StartGame(room roomModel.WaitingRoom) {
	// gameCfg := h.gameService.GameConfig(room.ID)
	// if err := h.gameService.CheckBeforeRegistration(*room, gameCfg); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	// mod, err := h.gameService.RegisterGame(gameCfg, room.PlayerIDs)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }
	// mod.OnPhaseChanged(func(mod game.Moderator) {
	// 	h.socketServer.BroadcastToRoom(
	// 		"/",
	// 		strconv.Itoa(int(mod.GameID())),
	// 		"phase_changed",
	// 		map[string]any{
	// 			"round":    mod.Scheduler().RoundID(),
	// 			"phase_id": mod.Scheduler().PhaseID(),
	// 		})
	// })
	// mod.StartGame()

	// // Store role assignment

	// set := make([]any, len(room.PlayerIDs)*2, len(room.PlayerIDs)*2)
	// for _, id := range room.PlayerIDs {
	// 	set = append(set, id, "in_game")
	// }
	// h.rdb.MSet(context.Background(), set...)

	// h.communicationService.BroadcastToRoom(room.ID, service.CommunicationEventMsg{
	// 	Event: "start_game",
	// 	Message: map[string]any{
	// 		"game_id": mod.GameID(),
	// 	},
	// })

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "Ok",
	// })
}

// // CheckBeforeRegistration checks the combination of room and game config before
// // registering a game. This makes sure the game runs properly without any unexpectation.
// func (gs gameService) CheckBeforeRegistration(room data.WaitingRoom, gameCfg model.GameSettings) error {
// 	if len(room.PlayerIDs) < int(gs.config.MinCapacity) {
// 		return fmt.Errorf("Invite more players to play!")
// 	} else if len(room.PlayerIDs) > int(gs.config.MaxCapacity) {
// 		return fmt.Errorf("Too many players!")
// 	}

// 	numberOfPlayers := len(room.PlayerIDs)
// 	if (numberOfPlayers%2 == 0 && numberOfPlayers/2 <= int(gameCfg.NumberWerewolves)) ||
// 		(numberOfPlayers%2 != 0 && numberOfPlayers/2 < int(gameCfg.NumberWerewolves)) {
// 		return fmt.Errorf("Unblanced number of werewolves!")
// 	}

// 	return nil
// }
