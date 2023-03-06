package role

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TwoSisterSuite struct {
	suite.Suite
	playerID types.PlayerID
}

func TestTwoSisterSuite(t *testing.T) {
	suite.Run(t, new(SeerSuite))
}

func (tss *TwoSisterSuite) SetupSuite() {
	tss.playerID = types.PlayerID("1")
}

func (tss TwoSisterSuite) TestNewTwoSister() {
	ctrl := gomock.NewController(tss.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	player := gamemock.NewMockPlayer(ctrl)

	game.EXPECT().Player(tss.playerID).Return(player).Times(1)

	ts, _ := NewTwoSister(game, tss.playerID)

	tss.Equal(vars.TwoSistersRoleID, ts.ID())
	tss.Equal(vars.NightPhaseID, ts.(*twoSister).phaseID)
	tss.Equal(vars.VillagerFactionID, ts.FactionID())
	tss.Equal(vars.FirstRound, ts.(*twoSister).beginRoundID)
	tss.Equal(player, ts.(*twoSister).player)
	tss.Equal(vars.Unlimited, ts.ActiveLimit(0))
	tss.Len(ts.(*twoSister).abilities, 1)
	tss.Equal(vars.PredictActionID, ts.(*twoSister).abilities[0].action.ID())
}

func (tss TwoSisterSuite) TestTwoSisterRegisterTurns() {
	ctrl := gomock.NewController(tss.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	player := gamemock.NewMockPlayer(ctrl)
	scheduler := gamemock.NewMockScheduler(ctrl)

	game.EXPECT().Player(tss.playerID).Return(player).Times(1)
	game.EXPECT().Scheduler().Return(scheduler).Times(1)
	player.EXPECT().ID().Return(tss.playerID).Times(1)

	ts, _ := NewTwoSister(game, tss.playerID)

	scheduler.EXPECT().AddPlayerTurn(&types.NewPlayerTurn{
		PhaseID:      ts.(*twoSister).phaseID,
		TurnID:       ts.(*twoSister).turnID,
		BeginRoundID: ts.(*twoSister).beginRoundID,
		PlayerID:     tss.playerID,
		RoleID:       ts.(*twoSister).id,
		ExpiredAfter: vars.One,
	})

	ts.RegisterTurn()
}
