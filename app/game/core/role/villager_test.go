package role

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/types"
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

	playerID := types.PlayerID("1")
	tests := []struct {
		name  string
		isNil bool
		err   string
		setup func()
	}{
		{
			name:  "Nil (Poll does not exist)",
			isNil: true,
			err:   "Poll does not exist ¯\\_(ツ)_/¯",
			setup: func() {
				mockGame.EXPECT().Poll(config.VillagerFactionID).Return(nil).Times(1)
			},
		},
		{
			name:  "Nil (Unable to join to the poll)",
			isNil: true,
			err:   "Unable to join to the poll ಠ_ಠ",
			setup: func() {
				mockGame.EXPECT().Poll(config.VillagerFactionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(false).Times(1)
			},
		},
		{
			name:  "Ok",
			isNil: false,
			setup: func() {
				mockGame.EXPECT().Poll(config.VillagerFactionID).Return(mockPoll).Times(1)
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

			if test.isNil {
				assert.Nil(t, villager)
				assert.NotNil(t, err)
				assert.Equal(t, test.err, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, config.VillagerRoleID, villager.ID())
				assert.Equal(t, config.DayPhaseID, villager.PhaseID())
				assert.Equal(t, config.VillagerFactionID, villager.FactionID())
				assert.Equal(t, config.VillagerTurnPriority, villager.Priority())
				assert.Equal(t, config.FirstRound, villager.BeginRound())
				assert.Equal(t, config.Unlimited, villager.ActiveLimit(config.VoteActionID))
			}
		})
	}
}
