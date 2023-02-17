package role

import (
	"testing"
	"uwwolf/game/enum"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TwoSisterSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	game     *gamemock.MockGame
	player   *gamemock.MockPlayer
	playerID enum.PlayerID
}

func TestTwoSisterSuite(t *testing.T) {
	suite.Run(t, new(TwoSisterSuite))
}

func (tss *TwoSisterSuite) SetupSuite() {
	tss.playerID = "1"
}

func (tss *TwoSisterSuite) SetupTest() {
	tss.ctrl = gomock.NewController(tss.T())
	tss.game = gamemock.NewMockGame(tss.ctrl)
	tss.player = gamemock.NewMockPlayer(tss.ctrl)
	tss.game.EXPECT().Player(tss.playerID).Return(tss.player).AnyTimes()
}

func (tss *TwoSisterSuite) TearDownTest() {
	tss.ctrl.Finish()
}

func (tss *TwoSisterSuite) TestNewTwoSister() {
	twoSister, _ := NewTwoSister(tss.game, tss.playerID)

	tss.Equal(enum.TwoSistersRoleID, twoSister.ID())
	tss.Equal(enum.NightPhaseID, twoSister.PhaseID())
	tss.Equal(enum.VillagerFactionID, twoSister.FactionID())
	tss.Equal(enum.TwoSistersTurnPriority, twoSister.Priority())
	tss.Equal(enum.FirstRound, twoSister.BeginRound())
	tss.Equal(enum.OneMore, twoSister.ActiveLimit(enum.RecognizeActionID))
}
