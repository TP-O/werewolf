package contract

import (
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type GameManger interface {
	// Game returns game instance by game ID.
	Game(gameID enum.GameID) Game

	// AddGame inserts new game instance to the game manager.
	// The old one can be overrided if it has the same game ID.
	AddGame(gameID enum.GameID, setting *types.GameSetting) Game
}
