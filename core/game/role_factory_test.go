package game

import (
	"fmt"
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewRole(t *testing.T) {
	playerID := types.PlayerID("1")

	tests := []struct {
		name           string
		expectedRoleID types.RoleID
		expectedErr    error
		setup          func(*gamemock.MockGame, *gamemock.MockPoll)
	}{
		{
			name:           "New Villager",
			expectedRoleID: vars.VillagerRoleID,
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.VillagerFactionID).Return(mp).Times(2)
				mp.EXPECT().AddElectors(playerID)
				mp.EXPECT().SetWeight(playerID, uint(1))
			},
		},
		{
			name:           "New Werewolf",
			expectedRoleID: vars.WerewolfRoleID,
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(mp).Times(2)
				mp.EXPECT().AddElectors(playerID)
				mp.EXPECT().SetWeight(playerID, uint(1))
			},
		},
		{
			name:           "New Hunter",
			expectedRoleID: vars.HunterRoleID,
			setup:          func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {},
		},
		{
			name:           "New Seer",
			expectedRoleID: vars.SeerRoleID,
			setup:          func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {},
		},
		{
			name:           "New TwoSister",
			expectedRoleID: vars.TwoSistersRoleID,
			setup:          func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {},
		},
		{
			name:        "Non-existent role",
			expectedErr: fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯"),
			setup:       func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			game := gamemock.NewMockGame(ctrl)
			player := gamemock.NewMockPlayer(ctrl)
			poll := gamemock.NewMockPoll(ctrl)

			game.EXPECT().Player(playerID).Return(player).AnyTimes()
			test.setup(game, poll)

			role, err := NewRole(test.expectedRoleID, game, playerID)

			if test.expectedErr == nil {
				assert.Equal(t, test.expectedRoleID, role.ID())
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}
