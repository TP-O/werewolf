package core

import (
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type manager struct {
	games map[types.GameId]contract.Game
}

var mangerInstance *manager

func NewManager() contract.GameManger {
	if mangerInstance == nil {
		mangerInstance = &manager{
			games: make(map[types.GameId]contract.Game),
		}
	}

	return mangerInstance
}

func (m *manager) Game(gameId types.GameId) contract.Game {
	return m.games[gameId]
}

func (m *manager) AddGame(gameId types.GameId, setting *types.GameSetting) contract.Game {
	m.games[gameId] = NewGame(gameId, setting)

	return m.Game(gameId)
}
