package service

import (
	"uwwolf/app/enum"
	"uwwolf/app/types"
	"uwwolf/db"
)

func ArePlayersReadyToPlay(playerIds ...types.PlayerId) bool {
	var numberOfReadyPlayers int64
	db.Client().
		Count(&numberOfReadyPlayers).
		Where(
			"id IN (?) AND statusId NOT NULL AND statusId <> ?",
			playerIds,
			enum.InGameStatus,
		)

	return len(playerIds) == int(numberOfReadyPlayers)
}
