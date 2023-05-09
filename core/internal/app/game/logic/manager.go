package game

import (
	"fmt"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/internal/config"
)

type manager struct {
	config          config.Game
	moderators      map[types.GameID]Moderator
	playerID2GameID map[types.PlayerID]types.GameID
}

type Manager interface {
	Moderator(gameID types.GameID) Moderator

	ModeratorOfPlayer(playerID types.PlayerID) Moderator

	// NewGame inserts new game instance to the game manager.
	RegisterGame(registration *types.GameRegistration) (Moderator, error)
}

var m Manager

func NewManager(config config.Game) Manager {
	if m == nil {
		m = &manager{
			config:          config,
			moderators:      make(map[types.GameID]Moderator),
			playerID2GameID: make(map[types.PlayerID]types.GameID),
		}
	}
	return m
}

func (m manager) Moderator(gameID types.GameID) Moderator {
	return m.moderators[gameID]
}

func (m manager) ModeratorOfPlayer(playerID types.PlayerID) Moderator {
	if gameID := m.playerID2GameID[playerID]; gameID.IsUnknown() {
		return nil
	} else {
		return m.moderators[gameID]
	}
}

func (m *manager) RegisterGame(registration *types.GameRegistration) (Moderator, error) {
	if m.moderators[registration.ID] != nil {
		return nil, fmt.Errorf("Game is already running!")
	}

	m.moderators[registration.ID] = NewModerator(m.config, registration)

	for _, playerID := range registration.GameInitialization.PlayerIDs {
		m.playerID2GameID[playerID] = registration.ID
	}

	return m.moderators[registration.ID], nil
}
