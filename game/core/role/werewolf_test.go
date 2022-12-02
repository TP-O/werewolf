package role

import (
	"testing"
	"uwwolf/game/enum"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewWerewolf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)
	mockPoll := gamemock.NewMockPoll(ctrl)

	playerID := enum.PlayerID("1")
	tests := []struct {
		name        string
		returnNil   bool
		expectedRrr string
		setup       func()
	}{
		{
			name:        "Failure (Poll does not exist)",
			returnNil:   true,
			expectedRrr: "Poll does not exist ¯\\_(ツ)_/¯",
			setup: func() {
				mockGame.EXPECT().Poll(enum.WerewolfFactionID).Return(nil).Times(1)
			},
		},
		{
			name:        "Failure (Unable to join to the poll)",
			returnNil:   true,
			expectedRrr: "Unable to join to the poll ಠ_ಠ",
			setup: func() {
				mockGame.EXPECT().Poll(enum.WerewolfFactionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(false).Times(1)
			},
		},
		{
			name:      "Ok",
			returnNil: false,
			setup: func() {
				mockGame.EXPECT().Poll(enum.WerewolfFactionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(true).Times(1)
				mockPoll.EXPECT().SetWeight(playerID, uint(1)).Times(1)
				mockGame.EXPECT().Player(playerID).Return(mockPlayer).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			werewolf, err := NewWerewolf(mockGame, playerID)

			if test.returnNil {
				assert.Nil(t, werewolf)
				assert.NotNil(t, err)
				assert.Equal(t, test.expectedRrr, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, enum.WerewolfRoleID, werewolf.ID())
				assert.Equal(t, enum.NightPhaseID, werewolf.PhaseID())
				assert.Equal(t, enum.WerewolfFactionID, werewolf.FactionID())
				assert.Equal(t, enum.WerewolfTurnPriority, werewolf.Priority())
				assert.Equal(t, enum.FirstRound, werewolf.BeginRound())
				assert.Equal(t, enum.Unlimited, werewolf.ActiveLimit(enum.VoteActionID))
			}
		})
	}
}
