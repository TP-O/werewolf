package contract

import "uwwolf/game/types"

type Manager interface {
	Moderator(gameID types.GameID) Moderator

	// NewGame inserts new game instance to the game manager.
	AddModerator(gameID types.GameID, moderator Moderator) (bool, error)
}
