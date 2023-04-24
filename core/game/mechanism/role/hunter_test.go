package role

// import (
// 	"testing"
// 	"uwwolf/game/types"
// 	"uwwolf/game/vars"
// 	mock_game "uwwolf/mock/game"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type HunterSuite struct {
// 	suite.Suite
// 	playerID types.PlayerID
// }

// func TestHunterSuite(t *testing.T) {
// 	suite.Run(t, new(HunterSuite))
// }

// func (hs *HunterSuite) SetupSuite() {
// 	hs.playerID = types.PlayerID("1")
// }

// func (hs HunterSuite) TestNewHunter() {
// 	ctrl := gomock.NewController(hs.T())
// 	defer ctrl.Finish()
// 	game := mock_game.NewMockGame(ctrl)
// 	player := mock_game.NewMockPlayer(ctrl)

// 	game.EXPECT().Player(hs.playerID).Return(player)

// 	h, _ := NewHunter(game, hs.playerID)

// 	hs.Equal(vars.HunterRoleID, h.ID())
// 	hs.Equal(vars.DayPhaseID, h.(*hunter).phaseID)
// 	hs.Equal(vars.VillagerFactionID, h.FactionID())
// 	hs.Equal(vars.FirstRound, h.(*hunter).beginRoundID)
// 	hs.Equal(player, h.(*hunter).player)
// 	hs.Equal(vars.OutOfTimes, h.ActiveTimes(0))
// 	hs.Len(h.(*hunter).abilities, 1)
// 	hs.Equal(vars.KillActionID, h.(*hunter).abilities[0].action.ID())
// }

// func (hs HunterSuite) TestOnAssign() {
// 	ctrl := gomock.NewController(hs.T())
// 	defer ctrl.Finish()
// 	game := mock_game.NewMockGame(ctrl)
// 	player := mock_game.NewMockPlayer(ctrl)

// 	game.EXPECT().Player(hs.playerID).Return(player)

// 	h, _ := NewHunter(game, hs.playerID)
// 	h.OnAssign()
// }

// func (hs HunterSuite) TestOnAfterDeath() {
// 	tests := []struct {
// 		name          string
// 		expectedLimit types.Times
// 		setup         func(*hunter, *mock_game.MockGame, *mock_game.MockScheduler)
// 	}{
// 		{
// 			name:          "Die at inactive phase",
// 			expectedLimit: vars.Once,
// 			setup: func(h *hunter, mg *mock_game.MockGame, ms *mock_game.MockScheduler) {
// 				mg.EXPECT().Scheduler().Return(ms).Times(4)
// 				ms.EXPECT().PhaseID().Return(vars.NightPhaseID)
// 				ms.EXPECT().RoundID().Return(vars.SecondRound).Times(2)
// 				ms.EXPECT().AddSlot(&types.NewTurnSlot{
// 					PhaseID:       h.phaseID,
// 					RoleID:        h.id,
// 					PlayedRoundID: vars.SecondRound,
// 					PlayerID:      hs.playerID,
// 					TurnID:        h.turnID,
// 				})
// 			},
// 		},
// 		{
// 			name:          "Die at active phase",
// 			expectedLimit: vars.Once,
// 			setup: func(h *hunter, mg *mock_game.MockGame, ms *mock_game.MockScheduler) {
// 				mg.EXPECT().Scheduler().Return(ms).Times(5)
// 				ms.EXPECT().PhaseID().Return(vars.DayPhaseID)
// 				ms.EXPECT().TurnID().Return(vars.MidTurn)
// 				ms.EXPECT().RoundID().Return(vars.SecondRound).Times(2)
// 				ms.EXPECT().AddSlot(&types.NewTurnSlot{
// 					PhaseID:       h.phaseID,
// 					RoleID:        h.id,
// 					PlayedRoundID: vars.SecondRound,
// 					PlayerID:      hs.playerID,
// 					TurnID:        vars.MidTurn + 1,
// 				})
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		hs.Run(test.name, func() {
// 			ctrl := gomock.NewController(hs.T())
// 			defer ctrl.Finish()
// 			game := mock_game.NewMockGame(ctrl)
// 			player := mock_game.NewMockPlayer(ctrl)
// 			scheduler := mock_game.NewMockScheduler(ctrl)

// 			game.EXPECT().Player(hs.playerID).Return(player)
// 			player.EXPECT().ID().Return(hs.playerID).AnyTimes()

// 			h, _ := NewHunter(game, hs.playerID)
// 			test.setup(h.(*hunter), game, scheduler)
// 			h.OnAfterDeath()

// 			hs.Equal(test.expectedLimit, h.ActiveTimes(0))
// 		})
// 	}
// }
