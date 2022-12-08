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

func Manager() contract.GameManger {
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
	if m.games[gameID] != nil {
		return nil
	}

	m.games[gameID] = NewGame(gameID, setting)

	return m.games[gameID]
}

func (m *manager) RemoveGame(gameID enum.GameID) bool {
	if m.games[gameID] == nil {
		return false
	}

	delete(m.games, gameID)

	return true
}
