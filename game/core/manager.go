package core

import (
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type manager struct {
	games map[enum.GameID]contract.Game
}

var mangerInstance *manager

func NewManager() contract.GameManger {
	if mangerInstance == nil {
		mangerInstance = &manager{
			games: make(map[enum.GameID]contract.Game),
		}
	}

	return mangerInstance
}

func (m *manager) Game(gameID enum.GameID) contract.Game {
	return m.games[gameID]
}

func (m *manager) AddGame(gameID enum.GameID, setting *types.GameSetting) contract.Game {
	m.games[gameID] = NewGame(gameID, setting)

	return m.games[gameID]
}
