package contract

import (
	"uwwolf/game/types"
)

type GameManger interface {
	// Game returns game instance by game ID.
	Game(gameID types.GameID) Game

	// NewGame inserts new game instance to the game manager.
	NewGame(setting *types.GameSetting) (Game, error)

	// Remove removes the game with the given ID.
	RemoveGame(gameID types.GameID) (bool, error)
}
