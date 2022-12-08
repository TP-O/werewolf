package contract

import (
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type GameManger interface {
	// Game returns game instance by game ID.
	Game(gameID enum.GameID) Game

	// AddGame inserts new game instance to the game manager.
	// Returns nil if the old one is existed.
	AddGame(gameID enum.GameID, setting *types.GameSetting) Game

	// Remove removes the game with the given ID.
	RemoveGame(gameID enum.GameID) bool
}
