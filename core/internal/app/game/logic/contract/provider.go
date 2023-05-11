package contract

import "uwwolf/internal/app/game/logic/types"

type Manager interface {
	Moderator(gameID types.GameID) Moderator

	ModeratorOfPlayer(playerID types.PlayerId) Moderator

	// NewGame inserts new game instance to the game manager.
	RegisterGame(registration *types.GameRegistration) (Moderator, error)
}
