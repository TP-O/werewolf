package game

import (
	"fmt"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

type manager struct {
	moderators map[types.GameID]contract.Moderator
}

var m contract.Manager

func Manager() contract.Manager {
	if m == nil {
		m = &manager{
			moderators: make(map[types.GameID]contract.Moderator),
		}
	}

	return m
}

func (m *manager) Moderator(gameID types.GameID) contract.Moderator {
	return m.moderators[gameID]
}

func (m *manager) AddModerator(gameID types.GameID, moderator contract.Moderator) (bool, error) {
	if m.moderators[gameID] != nil {
		return false, fmt.Errorf("Game is already running!")
	}

	m.moderators[gameID] = moderator
	moderator.SetGameID(gameID)

	return true, nil
}
