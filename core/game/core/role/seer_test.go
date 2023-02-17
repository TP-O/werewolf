package role

import (
	"testing"
	"uwwolf/game/enum"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type SeerSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	game     *gamemock.MockGame
	player   *gamemock.MockPlayer
	playerID enum.PlayerID
}

func TestSeerSuite(t *testing.T) {
	suite.Run(t, new(SeerSuite))
}

func (ss *SeerSuite) SetupSuite() {
	ss.playerID = "1"
}

func (ss *SeerSuite) SetupTest() {
	ss.ctrl = gomock.NewController(ss.T())
	ss.game = gamemock.NewMockGame(ss.ctrl)
	ss.player = gamemock.NewMockPlayer(ss.ctrl)
	ss.game.EXPECT().Player(ss.playerID).Return(ss.player).AnyTimes()
}

func (ss *SeerSuite) TearDownTest() {
	ss.ctrl.Finish()
}

func (ss *SeerSuite) TestNewSeer() {
	seer, _ := NewSeer(ss.game, ss.playerID)

	ss.Equal(enum.SeerRoleID, seer.ID())
	ss.Equal(enum.NightPhaseID, seer.PhaseID())
	ss.Equal(enum.VillagerFactionID, seer.FactionID())
	ss.Equal(enum.SeerTurnPriority, seer.Priority())
	ss.Equal(enum.Round(2), seer.BeginRound())
	ss.Equal(enum.Unlimited, seer.ActiveLimit(enum.PredictActionID))
}
