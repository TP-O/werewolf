package role

import (
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type HunterSuite struct {
	suite.Suite
	playerId types.PlayerId
}

func TestHunterSuite(t *testing.T) {
	suite.Run(t, new(HunterSuite))
}

func (hs *HunterSuite) SetupSuite() {
	hs.playerId = types.PlayerId("1")
}

func (hs HunterSuite) TestNewHunter() {
	ctrl := gomock.NewController(hs.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	moderator.EXPECT().World().Return(nil)

	h, _ := NewHunter(moderator, hs.playerId)

	hs.Equal(constants.HunterRoleId, h.Id())
	hs.Equal(constants.DayPhaseId, h.(*hunter).phaseId)
	hs.Equal(constants.VillagerFactionId, h.FactionId())
	hs.Equal(constants.FirstRound, h.(*hunter).beginRound)
	hs.Equal(hs.playerId, h.(*hunter).playerId)
	hs.Equal(constants.OutOfTimes, h.ActiveTimes(0))
	hs.Len(h.(*hunter).abilities, 1)
	hs.Equal(constants.KillActionId, h.(*hunter).abilities[0].action.Id())
	hs.True(h.(*hunter).abilities[0].isImmediate)
}

func (hs HunterSuite) TestOnAfterAssign() {
	ctrl := gomock.NewController(hs.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	world := mock_game_logic.NewMockWorld(ctrl)

	moderator.EXPECT().World().Return(world)

	h, _ := NewHunter(moderator, hs.playerId)
	h.OnAfterAssign()
}

func (hs HunterSuite) TestOnAfterDeath() {
	tests := []struct {
		name          string
		expectedLimit types.Times
		setup         func(*hunter, *mock_game_logic.MockModerator, *mock_game_logic.MockScheduler)
	}{
		{
			name:          "Die too soon",
			expectedLimit: constants.OutOfTimes,
			setup: func(h *hunter, mm *mock_game_logic.MockModerator, ms *mock_game_logic.MockScheduler) {
				mm.EXPECT().Scheduler().Return(ms).Times(2)
				ms.EXPECT().PhaseId().Return(constants.NightPhaseId)
				ms.EXPECT().Round().Return(constants.ZeroRound)
			},
		},
		{
			name:          "Die at inactive phase",
			expectedLimit: constants.Once,
			setup: func(h *hunter, mm *mock_game_logic.MockModerator, ms *mock_game_logic.MockScheduler) {
				mm.EXPECT().Scheduler().Return(ms).Times(4)
				ms.EXPECT().PhaseId().Return(constants.NightPhaseId)
				ms.EXPECT().Round().Return(constants.SecondRound).Times(2)
				ms.EXPECT().AddSlot(types.AddTurnSlot{
					PhaseId:  h.phaseId,
					PlayerId: hs.playerId,
					Turn:     h.turn,
					TurnSlot: types.TurnSlot{
						RoleId:      h.id,
						PlayedRound: constants.SecondRound,
					},
				})
			},
		},
		{
			name:          "Die at active phase",
			expectedLimit: constants.Once,
			setup: func(h *hunter, mm *mock_game_logic.MockModerator, ms *mock_game_logic.MockScheduler) {
				mm.EXPECT().Scheduler().Return(ms).Times(5)
				ms.EXPECT().PhaseId().Return(constants.DayPhaseId)
				ms.EXPECT().Turn().Return(constants.MidTurn)
				ms.EXPECT().Round().Return(constants.SecondRound).Times(2)
				ms.EXPECT().AddSlot(types.AddTurnSlot{
					PhaseId:  h.phaseId,
					PlayerId: hs.playerId,
					Turn:     constants.MidTurn + 1,
					TurnSlot: types.TurnSlot{
						RoleId:      h.id,
						PlayedRound: constants.SecondRound,
					},
				})
			},
		},
	}

	for _, test := range tests {
		hs.Run(test.name, func() {
			ctrl := gomock.NewController(hs.T())
			defer ctrl.Finish()
			moderator := mock_game_logic.NewMockModerator(ctrl)
			scheduler := mock_game_logic.NewMockScheduler(ctrl)

			moderator.EXPECT().World().Return(nil)

			h, _ := NewHunter(moderator, hs.playerId)
			test.setup(h.(*hunter), moderator, scheduler)
			h.OnAfterDeath()

			hs.Equal(test.expectedLimit, h.ActiveTimes(0))
		})
	}
}
