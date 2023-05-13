package role

import (
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TwoSisterSuite struct {
	suite.Suite
	playerId types.PlayerId
}

func TestTwoSisterSuite(t *testing.T) {
	suite.Run(t, new(TwoSisterSuite))
}

func (tss *TwoSisterSuite) SetupSuite() {
	tss.playerId = types.PlayerId("1")
}

func (tss TwoSisterSuite) TestNewTwoSister() {
	ctrl := gomock.NewController(tss.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)

	moderator.EXPECT().World().Return(nil)
	moderator.EXPECT().RegisterActionExecution(gomock.Any())

	ts, _ := NewTwoSister(moderator, tss.playerId)

	tss.Equal(constants.TwoSistersRoleId, ts.Id())
	tss.Equal(constants.NightPhaseId, ts.(*twoSister).phaseId)
	tss.Equal(constants.VillagerFactionId, ts.FactionId())
	tss.Equal(constants.FirstRound, ts.(*twoSister).beginRound)
	tss.Equal(tss.playerId, ts.(*twoSister).playerId)
	tss.Equal(constants.OutOfTimes, ts.ActiveTimes(0))
	tss.Len(ts.(*twoSister).abilities, 0)
}
