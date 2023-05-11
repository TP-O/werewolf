package role

// import (
// 	"testing"
// 	"uwwolf/internal/app/game/logic/types"
// 	"uwwolf/game/vars"
// 	mock_game "uwwolf/mock/game"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type TwoSisterSuite struct {
// 	suite.Suite
// 	playerID types.PlayerId
// }

// func TestTwoSisterSuite(t *testing.T) {
// 	suite.Run(t, new(SeerSuite))
// }

// func (tss *TwoSisterSuite) SetupSuite() {
// 	tss.playerID = types.PlayerId("1")
// }

// func (tss TwoSisterSuite) TestNewTwoSister() {
// 	ctrl := gomock.NewController(tss.T())
// 	defer ctrl.Finish()
// 	game := mock_game.NewMockGame(ctrl)
// 	player := mock_game.NewMockPlayer(ctrl)

// 	game.EXPECT().Player(tss.playerID).Return(player)

// 	ts, _ := NewTwoSister(game, tss.playerID)

// 	tss.Equal(vars.TwoSistersRoleID, ts.Id())
// 	tss.Equal(vars.NightPhaseID, ts.(*twoSister).phaseID)
// 	tss.Equal(vars.VillagerFactionID, ts.FactionID())
// 	tss.Equal(vars.FirstRound, ts.(*twoSister).beginRoundID)
// 	tss.Equal(player, ts.(*twoSister).player)
// 	tss.Equal(vars.UnlimitedTimes, ts.ActiveTimes(0))
// 	tss.Len(ts.(*twoSister).abilities, 1)
// 	tss.Equal(vars.PredictActionID, ts.(*twoSister).abilities[0].action.Id())
// }
