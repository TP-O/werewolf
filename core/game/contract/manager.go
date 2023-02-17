package contract

import (
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type GameManger interface {
	// Game returns game instance by game ID.
	Game(gameID enum.GameID) Game

	// NewGame inserts new game instance to the game manager.
	NewGame(setting *types.GameSetting) (Game, error)

	// Remove removes the game with the given ID.
	RemoveGame(gameID enum.GameID) (bool, error)
}
