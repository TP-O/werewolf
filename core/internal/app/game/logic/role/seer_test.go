package role

import (
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type SeerSuite struct {
	suite.Suite
	playerId types.PlayerId
}

func TestSeerSuite(t *testing.T) {
	suite.Run(t, new(SeerSuite))
}

func (ss *SeerSuite) SetupSuite() {
	ss.playerId = types.PlayerId("1")
}

func (ss SeerSuite) TestNewSeer() {
	ctrl := gomock.NewController(ss.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	moderator.EXPECT().World().Return(nil)

	s, _ := NewSeer(moderator, ss.playerId)

	ss.Equal(constants.SeerRoleId, s.Id())
	ss.Equal(constants.NightPhaseId, s.(*seer).phaseId)
	ss.Equal(constants.VillagerFactionId, s.FactionId())
	ss.Equal(constants.SecondRound, s.(*seer).beginRound)
	ss.Equal(ss.playerId, s.(*seer).playerId)
	ss.Equal(constants.UnlimitedTimes, s.ActiveTimes(0))
	ss.Len(s.(*seer).abilities, 1)
	ss.Equal(constants.PredictActionId, s.(*seer).abilities[0].action.Id())
	ss.True(s.(*seer).abilities[0].isImmediate)
}
