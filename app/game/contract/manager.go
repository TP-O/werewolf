package contract

import (
	"uwwolf/app/types"
)

type GameManger interface {
	// Game returns game instance by game id.
	Game(gameId types.GameId) Game

	// AddGame inserts new game instance to game manager.
	// Old instance can be overrided if it has the same game id.
	AddGame(gameId types.GameId, setting *types.GameSetting) Game
}
