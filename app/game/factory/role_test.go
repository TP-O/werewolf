package factory

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockPlayer := gamemock.NewMockPlayer(ctrl)
	mockPoll := gamemock.NewMockPoll(ctrl)

	playerID := types.PlayerID("1")
	mockGame.EXPECT().Player(playerID).Return(mockPlayer).AnyTimes()

	tests := []struct {
		name           string
		expectedRoleID types.RoleID
		expectedErr    string
		setup          func()
	}{
		{
			name:           "New Villager",
			expectedRoleID: config.VillagerRoleID,
			setup: func() {
				mockGame.EXPECT().Poll(config.VillagerFactionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(true).Times(1)
				mockPoll.EXPECT().SetWeight(playerID, uint(1)).Times(1)
			},
		},
		{
			name:           "New Werewolf",
			expectedRoleID: config.WerewolfRoleID,
			setup: func() {
				mockGame.EXPECT().Poll(config.WerewolfFactionID).Return(mockPoll).Times(1)
				mockPoll.EXPECT().AddElectors(playerID).Return(true).Times(1)
				mockPoll.EXPECT().SetWeight(playerID, uint(1)).Times(1)
			},
		},
		{
			name:           "New Hunter",
			expectedRoleID: config.HunterRoleID,
			setup:          func() {},
		},
		{
			name:           "New Seer",
			expectedRoleID: config.SeerRoleID,
			setup:          func() {},
		},
		{
			name:           "New TwoSister",
			expectedRoleID: config.TwoSistersRoleID,
			setup:          func() {},
		},
		{
			name:        "Non-existent role",
			expectedErr: "Non-existent role ¯\\_ಠ_ಠ_/¯",
			setup:       func() {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			role, err := NewRole(test.expectedRoleID, mockGame, playerID)

			if test.expectedErr == "" {
				assert.Equal(t, test.expectedRoleID, role.ID())
			} else {
				assert.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}
