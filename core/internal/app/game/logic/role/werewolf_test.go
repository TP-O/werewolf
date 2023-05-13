package role

import (
	"errors"
	"testing"
	"uwwolf/internal/app/game/logic/action"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type WerewolfSuite struct {
	suite.Suite
	playerId types.PlayerId
}

func TestWerewolfSuite(t *testing.T) {
	suite.Run(t, new(WerewolfSuite))
}

func (ws *WerewolfSuite) SetupSuite() {
	ws.playerId = types.PlayerId("1")
}

func (ws WerewolfSuite) TestNewWerewolf() {
	tests := []struct {
		name        string
		expectedErr error
		setup       func(*mock_game_logic.MockModerator, *mock_game_logic.MockWorld, *mock_game_logic.MockPoll)
	}{
		{
			name:        "Failure (Poll does not exist)",
			expectedErr: errors.New("Poll does not exist ¯\\_(ツ)_/¯"),
			setup: func(
				mm *mock_game_logic.MockModerator,
				mw *mock_game_logic.MockWorld,
				mp *mock_game_logic.MockPoll) {
				mw.EXPECT().Poll(constants.WerewolfFactionId).Return(nil)
			},
		},
		{
			name: "Ok",
			setup: func(
				mm *mock_game_logic.MockModerator,
				mw *mock_game_logic.MockWorld,
				mp *mock_game_logic.MockPoll) {
				mw.EXPECT().Poll(constants.WerewolfFactionId).Return(mp).Times(2)
				mp.EXPECT().AddElectors(ws.playerId)
				mp.EXPECT().SetWeight(ws.playerId, uint(1))
			},
		},
	}

	for _, test := range tests {
		ws.Run(test.name, func() {
			ctrl := gomock.NewController(ws.T())
			defer ctrl.Finish()
			moderator := mock_game_logic.NewMockModerator(ctrl)
			world := mock_game_logic.NewMockWorld(ctrl)
			poll := mock_game_logic.NewMockPoll(ctrl)

			moderator.EXPECT().World().Return(world)

			test.setup(moderator, world, poll)

			w, err := NewWerewolf(moderator, ws.playerId)

			if test.expectedErr != nil {
				ws.Nil(w)
				ws.NotNil(err)
				ws.Equal(test.expectedErr, err)
			} else {
				ws.Nil(err)
				ws.Equal(constants.WerewolfRoleId, w.Id())
				ws.Equal(constants.NightPhaseId, w.(*werewolf).phaseId)
				ws.Equal(constants.WerewolfFactionId, w.FactionId())
				ws.Equal(constants.FirstRound, w.(*werewolf).beginRound)
				ws.Equal(ws.playerId, w.(*werewolf).playerId)
				ws.Equal(constants.UnlimitedTimes, w.ActiveTimes(0))
				ws.Len(w.(*werewolf).abilities, 1)
				ws.Equal(action.VoteActionId, w.(*werewolf).abilities[0].action.Id())
				ws.True(w.(*werewolf).abilities[0].isImmediate)
			}
		})
	}
}

func (ws WerewolfSuite) TestAfterOnAssign() {
	ctrl := gomock.NewController(ws.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	scheduler := mock_game_logic.NewMockScheduler(ctrl)
	world := mock_game_logic.NewMockWorld(ctrl)
	poll := mock_game_logic.NewMockPoll(ctrl)

	// Mock for New fuction
	moderator.EXPECT().World().Return(world)
	moderator.EXPECT().Scheduler().Return(scheduler)
	world.EXPECT().Poll(constants.WerewolfFactionId).Return(poll).Times(2)
	poll.EXPECT().AddElectors(ws.playerId)
	poll.EXPECT().SetWeight(ws.playerId, uint(1))

	w, _ := NewWerewolf(moderator, ws.playerId)

	moderator.EXPECT().World().Return(world)
	world.EXPECT().Poll(constants.WerewolfFactionId).Return(poll)
	poll.EXPECT().AddElectors(ws.playerId)
	scheduler.EXPECT().AddSlot(types.NewTurnSlot{
		PhaseId:    w.(*werewolf).phaseId,
		Turn:       w.(*werewolf).turn,
		BeginRound: w.(*werewolf).beginRound,
		PlayerId:   ws.playerId,
		RoleId:     w.Id(),
	})

	w.OnAfterAssign()
}

func (ws WerewolfSuite) TestOnAfterRevoke() {
	ctrl := gomock.NewController(ws.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	scheduler := mock_game_logic.NewMockScheduler(ctrl)
	world := mock_game_logic.NewMockWorld(ctrl)
	poll := mock_game_logic.NewMockPoll(ctrl)

	// Mock for New fuction
	moderator.EXPECT().World().Return(world)
	world.EXPECT().Poll(constants.WerewolfFactionId).Return(poll).Times(2)
	poll.EXPECT().AddElectors(ws.playerId)
	poll.EXPECT().SetWeight(ws.playerId, uint(1))

	w, _ := NewWerewolf(moderator, ws.playerId)

	moderator.EXPECT().World().Return(world)
	moderator.EXPECT().Scheduler().Return(scheduler)
	world.EXPECT().Poll(constants.WerewolfFactionId).Return(poll)
	poll.EXPECT().RemoveElector(ws.playerId)
	scheduler.EXPECT().RemoveSlot(types.RemovedTurnSlot{
		PhaseId:  w.(*werewolf).phaseId,
		PlayerId: ws.playerId,
		RoleId:   w.Id(),
	})

	w.OnAfterRevoke()
}
