package role

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type VillagerSuite struct {
	suite.Suite
	playerID types.PlayerID
}

func TestVillagerSuite(t *testing.T) {
	suite.Run(t, new(VillagerSuite))
}

func (vs *VillagerSuite) SetupSuite() {
	vs.playerID = types.PlayerID("1")
}

func (vs VillagerSuite) TestNewVillager() {
	tests := []struct {
		name        string
		expectedErr string
		setup       func(*gamemock.MockGame, *gamemock.MockPoll)
	}{
		{
			name:        "Failure (Poll does not exist)",
			expectedErr: "Poll does not exist ¯\\_(ツ)_/¯",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.VillagerFactionID).Return(nil).Times(1)
			},
		},
		{
			name: "Ok",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.VillagerFactionID).Return(mp).Times(2)
				mp.EXPECT().AddElectors(vs.playerID).Times(1)
				mp.EXPECT().SetWeight(vs.playerID, uint(1)).Times(1)
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			game := gamemock.NewMockGame(ctrl)
			player := gamemock.NewMockPlayer(ctrl)
			poll := gamemock.NewMockPoll(ctrl)

			game.EXPECT().Player(vs.playerID).Return(player).AnyTimes()
			test.setup(game, poll)

			v, err := NewVillager(game, vs.playerID)

			if test.expectedErr != "" {
				vs.Nil(v)
				vs.NotNil(err)
				vs.Equal(test.expectedErr, err.Error())
			} else {
				vs.Nil(err)
				vs.Equal(vars.VillagerRoleID, v.ID())
				vs.Equal(vars.DayPhaseID, v.(*villager).phaseID)
				vs.Equal(vars.VillagerFactionID, v.FactionID())
				vs.Equal(vars.FirstRound, v.(*villager).beginRoundID)
				vs.Equal(player, v.(*villager).player)
				vs.Equal(vars.UnlimitedTimes, v.ActiveTimes(0))
				vs.Len(v.(*villager).abilities, 1)
				vs.Equal(vars.VoteActionID, v.(*villager).abilities[0].action.ID())
			}
		})
	}
}
