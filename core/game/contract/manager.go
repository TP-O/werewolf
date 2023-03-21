package contract

import "uwwolf/game/types"

type Manager interface {
	// NewGame inserts new game instance to the game manager.
	AddModerator(gameID types.GameID, moderator Moderator) (bool, error)
}
