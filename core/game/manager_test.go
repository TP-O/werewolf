package game

import (
	"fmt"
	"testing"
	"uwwolf/config"
	"uwwolf/game/types"
	mock_game "uwwolf/mock/game"

	"github.com/stretchr/testify/suite"
)

type ManagerSuite struct {
	suite.Suite
	gameID types.GameID
}

func (ms *ManagerSuite) SetupSuite() {
	ms.gameID = types.GameID(1)
}

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (ms ManagerSuite) TestManager() {
	m1 := Manager(config.Game{})
	m2 := Manager(config.Game{})

	ms.Require().NotNil(m1)
	ms.Require().True(m1 == m2)
	ms.NotNil(m1.(*manager).moderators)
}

func (ms ManagerSuite) TestModerator() {
	mod := mock_game.NewMockModerator(nil)

	m := Manager(config.Game{})
	m.(*manager).moderators[ms.gameID] = mod

	ms.Require().True(m.Moderator(ms.gameID) == mod)
}

func (ms ManagerSuite) TestRegiserGame() {
	tests := []struct {
		name        string
		expectedErr error
		setup       func(m *manager)
	}{
		{
			name:        "Failure (GameID is ready registered)",
			expectedErr: fmt.Errorf("Game is already running!"),
			setup: func(m *manager) {
				m.moderators[ms.gameID] = mock_game.NewMockModerator(nil)
			},
		},
		{
			name: "Ok",
			setup: func(m *manager) {
				m.moderators[ms.gameID] = nil
			},
		},
	}

	for _, test := range tests {
		ms.Run(test.name, func() {
			m := Manager(config.Game{})
			test.setup(m.(*manager))

			mod, err := m.RegisterGame(&types.GameRegistration{
				ID: ms.gameID,
			})

			ms.Require().Equal(test.expectedErr, err)
			if test.expectedErr == nil {
				ms.Require().NotNil(mod)
			}
		})
	}
}
