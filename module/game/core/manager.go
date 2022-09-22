package core

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

type manager struct {
	games map[types.GameId]contract.Game
}

var mangerInstance *manager

func New() contract.GameManger {
	if mangerInstance == nil {
		mangerInstance = &manager{}
	}

	return mangerInstance
}

func (m *manager) Game(gameId types.GameId) contract.Game {
	return m.games[gameId]
}

func (m *manager) AddGame(setting *types.GameSetting) {
	m.games[setting.Id] = NewGame(setting)
}
