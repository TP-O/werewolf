package role

// import (
// 	"fmt"
// 	"testing"
// 	"uwwolf/internal/app/game/logic/types"
// 	"uwwolf/game/vars"
// 	mock_game "uwwolf/mock/game"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestNewRole(t *testing.T) {
// 	playerID := types.PlayerID("1")

// 	tests := []struct {
// 		name           string
// 		expectedRoleID types.RoleID
// 		expectedErr    error
// 		setup          func(*mock_game.MockGame, *mock_game.MockPoll)
// 	}{
// 		{
// 			name:           "New Villager",
// 			expectedRoleID: vars.VillagerRoleID,
// 			setup: func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {
// 				mg.EXPECT().Poll(vars.VillagerFactionID).Return(mp).Times(2)
// 				mp.EXPECT().AddElectors(playerID)
// 				mp.EXPECT().SetWeight(playerID, uint(1))
// 			},
// 		},
// 		{
// 			name:           "New Werewolf",
// 			expectedRoleID: vars.WerewolfRoleID,
// 			setup: func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {
// 				mg.EXPECT().Poll(vars.WerewolfFactionID).Return(mp).Times(2)
// 				mp.EXPECT().AddElectors(playerID)
// 				mp.EXPECT().SetWeight(playerID, uint(1))
// 			},
// 		},
// 		{
// 			name:           "New Hunter",
// 			expectedRoleID: vars.HunterRoleID,
// 			setup:          func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {},
// 		},
// 		{
// 			name:           "New Seer",
// 			expectedRoleID: vars.SeerRoleID,
// 			setup:          func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {},
// 		},
// 		{
// 			name:           "New TwoSister",
// 			expectedRoleID: vars.TwoSistersRoleID,
// 			setup:          func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {},
// 		},
// 		{
// 			name:        "Non-existent role",
// 			expectedErr: fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯"),
// 			setup:       func(mg *mock_game.MockGame, mp *mock_game.MockPoll) {},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)
// 			player := mock_game.NewMockPlayer(ctrl)
// 			poll := mock_game.NewMockPoll(ctrl)

// 			game.EXPECT().Player(playerID).Return(player).AnyTimes()
// 			test.setup(game, poll)

// 			role, err := NewRole(test.expectedRoleID, game, playerID)

// 			if test.expectedErr == nil {
// 				assert.Equal(t, test.expectedRoleID, role.ID())
// 			} else {
// 				assert.Equal(t, test.expectedErr, err)
// 			}
// 		})
// 	}
// }
