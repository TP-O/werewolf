package core

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

type gameManager struct {
	games map[types.GameId]contract.Game
}

var gameMangerInstance *gameManager

func New() contract.GameManger {
	if gameMangerInstance == nil {
		gameMangerInstance = &gameManager{}
	}

	return gameMangerInstance
}

func (gm *gameManager) Game(gameId types.GameId) contract.Game {
	return gm.games[gameId]
}

func (gm *gameManager) AddGame(setting *types.GameSetting) {
	gm.games[setting.Id] = NewGame(setting)
}
