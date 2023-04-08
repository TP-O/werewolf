package contract

import "uwwolf/game/types"

type Manager interface {
	Moderator(gameID types.GameID) Moderator

	ModeratorOfPlayer(playerID types.PlayerID) Moderator

	// NewGame inserts new game instance to the game manager.
	RegisterGame(registration *types.GameRegistration) (Moderator, error)
}
