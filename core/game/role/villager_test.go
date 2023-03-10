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
				mg.EXPECT().Poll(vars.VillagerFactionID).Return(nil)
			},
		},
		{
			name: "Ok",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPoll) {
				mg.EXPECT().Poll(vars.VillagerFactionID).Return(mp).Times(2)
				mp.EXPECT().AddElectors(vs.playerID)
				mp.EXPECT().SetWeight(vs.playerID, uint(1))
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

func (vs VillagerSuite) TestOnRevoke() {
	ctrl := gomock.NewController(vs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	player := gamemock.NewMockPlayer(ctrl)
	scheduler := gamemock.NewMockScheduler(ctrl)
	poll := gamemock.NewMockPoll(ctrl)

	// Mock for New fuction
	game.EXPECT().Scheduler().Return(scheduler)
	game.EXPECT().Player(vs.playerID).Return(player)
	game.EXPECT().Poll(vars.VillagerFactionID).Return(poll).Times(2)
	poll.EXPECT().AddElectors(vs.playerID)
	poll.EXPECT().SetWeight(vs.playerID, uint(1))

	game.EXPECT().Poll(vars.VillagerFactionID).Return(poll).Times(2)
	game.EXPECT().Poll(vars.WerewolfFactionID).Return(poll)
	poll.EXPECT().RemoveElector(vs.playerID)
	poll.EXPECT().RemoveCandidate(vs.playerID).Times(2)
	player.EXPECT().ID().Return(vs.playerID).Times(4)

	v, _ := NewVillager(game, vs.playerID)

	scheduler.EXPECT().RemoveSlot(&types.RemovedTurnSlot{
		PhaseID:  v.(*villager).phaseID,
		PlayerID: vs.playerID,
		RoleID:   v.ID(),
	})

	v.OnRevoke()
}

func (vs VillagerSuite) TestOnAssign() {
	ctrl := gomock.NewController(vs.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	player := gamemock.NewMockPlayer(ctrl)
	scheduler := gamemock.NewMockScheduler(ctrl)
	poll := gamemock.NewMockPoll(ctrl)

	// Mock for New fuction
	game.EXPECT().Scheduler().Return(scheduler)
	game.EXPECT().Player(vs.playerID).Return(player)
	game.EXPECT().Poll(vars.VillagerFactionID).Return(poll).Times(2)
	poll.EXPECT().AddElectors(vs.playerID)
	poll.EXPECT().SetWeight(vs.playerID, uint(1))

	game.EXPECT().Poll(vars.VillagerFactionID).Return(poll)
	game.EXPECT().Poll(vars.WerewolfFactionID).Return(poll)
	poll.EXPECT().AddCandidates(vs.playerID).Times(2)
	player.EXPECT().ID().Return(vs.playerID).Times(3)

	v, _ := NewVillager(game, vs.playerID)

	scheduler.EXPECT().AddSlot(&types.NewTurnSlot{
		PhaseID:      v.(*villager).phaseID,
		TurnID:       v.(*villager).turnID,
		BeginRoundID: v.(*villager).beginRoundID,
		PlayerID:     vs.playerID,
		RoleID:       v.ID(),
	})

	v.OnAssign()
}
