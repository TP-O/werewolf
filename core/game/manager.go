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

func (m *manager) RegisterGame(registration *types.GameRegistration) (contract.Moderator, error) {
	if m.moderators[registration.ID] != nil {
		return nil, fmt.Errorf("Game is already running!")
	}

	m.moderators[registration.ID] = NewModerator(registration)
	return m.moderators[registration.ID], nil
}
