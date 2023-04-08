package game

import (
	"fmt"
	"uwwolf/config"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

type manager struct {
	config          config.Game
	moderators      map[types.GameID]contract.Moderator
	playerID2GameID map[types.PlayerID]types.GameID
}

var m contract.Manager

func Manager(config config.Game) contract.Manager {
	if m == nil {
		m = &manager{
			config:          config,
			moderators:      make(map[types.GameID]contract.Moderator),
			playerID2GameID: make(map[types.PlayerID]types.GameID),
		}
	}
	return m
}

func (m manager) Moderator(gameID types.GameID) contract.Moderator {
	return m.moderators[gameID]
}

func (m manager) ModeratorOfPlayer(playerID types.PlayerID) contract.Moderator {
	if gameID := m.playerID2GameID[playerID]; gameID.IsUnknown() {
		return nil
	} else {
		return m.moderators[gameID]
	}
}

func (m *manager) RegisterGame(registration *types.GameRegistration) (contract.Moderator, error) {
	if m.moderators[registration.ID] != nil {
		return nil, fmt.Errorf("Game is already running!")
	}

	m.moderators[registration.ID] = NewModerator(m.config, registration)

	for _, playerID := range registration.GameInitialization.PlayerIDs {
		m.playerID2GameID[playerID] = registration.ID
	}

	return m.moderators[registration.ID], nil
}
