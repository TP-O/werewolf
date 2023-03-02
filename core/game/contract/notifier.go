package contract

import "uwwolf/game/types"

// Notifier sends notification or game state changes.
type Notifier interface {
	// NotifyPlayer sends message to the specific player.
	NotifyPlayer(playerID types.PlayerID, msg string)

	// NotifyPlayers sends message to the  group of players.
	NotifyPlayers(playerIDs []types.PlayerID, msg string)

	// NotifyGame sends message to all players in the specific game.
	NotifyGame(msg string, gameID types.GameID)
}
