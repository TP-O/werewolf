package role

import (
	"testing"
	"uwwolf/game/enum"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewVillager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)
	mockPoll := gamemock.NewMockPoll(ctrl)

	playerID := enum.PlayerID("1")
	tests := []struct {
		name        string
		expectedErr string
		returnNil   bool
		setup       func()
	}{
		{
			name:        "Failure (Poll does not exist)",
			returnNil:   true,
			expectedErr: "Poll does not exist ¯\\_(ツ)_/¯",
			setup: func() {
				mockGame.EXPECT().Poll(enum.VillagerFactionID).Return(nil).Times(1)
			},
		},
		{
			name:        "Failure (Unable to join to the poll)",
			returnNil:   true,
			expectedErr: "Unable to join to the poll ಠ_ಠ",
			setup: func() {
				mockGame.EXPECT().Poll(enum.VillagerFactionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(false).Times(1)
			},
		},
		{
			name:      "Ok",
			returnNil: false,
			setup: func() {
				mockGame.EXPECT().Poll(enum.VillagerFactionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(true).Times(1)
				mockPoll.EXPECT().SetWeight(playerID, uint(1)).Times(1)
				mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			villager, err := NewVillager(mockGame, playerID)

			if test.returnNil {
				assert.Nil(t, villager)
				assert.NotNil(t, err)
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, enum.VillagerRoleID, villager.ID())
				assert.Equal(t, enum.DayPhaseID, villager.PhaseID())
				assert.Equal(t, enum.VillagerFactionID, villager.FactionID())
				assert.Equal(t, enum.VillagerTurnPriority, villager.Priority())
				assert.Equal(t, enum.FirstRound, villager.BeginRound())
				assert.Equal(t, enum.Unlimited, villager.ActiveLimit(enum.VoteActionID))
			}
		})
	}
}
