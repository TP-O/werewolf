package role

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type SeerSuite struct {
	suite.Suite
	playerID types.PlayerID
}

func TestSeerSuite(t *testing.T) {
	suite.Run(t, new(SeerSuite))
}

func (ss *SeerSuite) SetupSuite() {
	ss.playerID = types.PlayerID("1")
}

func (ss SeerSuite) TestNewSeer() {
	ctrl := gomock.NewController(ss.T())
	defer ctrl.Finish()
	game := gamemock.NewMockGame(ctrl)
	player := gamemock.NewMockPlayer(ctrl)

	game.EXPECT().Player(ss.playerID).Return(player)

	s, _ := NewSeer(game, ss.playerID)

	ss.Equal(vars.SeerRoleID, s.ID())
	ss.Equal(vars.NightPhaseID, s.(*seer).phaseID)
	ss.Equal(vars.VillagerFactionID, s.FactionID())
	ss.Equal(vars.SecondRound, s.(*seer).beginRoundID)
	ss.Equal(player, s.(*seer).player)
	ss.Equal(vars.UnlimitedTimes, s.ActiveTimes(0))
	ss.Len(s.(*seer).abilities, 1)
	ss.Equal(vars.PredictActionID, s.(*seer).abilities[0].action.ID())
}
