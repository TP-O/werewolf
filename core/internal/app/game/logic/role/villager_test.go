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

type VillagerSuite struct {
	suite.Suite
	playerId types.PlayerId
}

func TestVillagerSuite(t *testing.T) {
	suite.Run(t, new(VillagerSuite))
}

func (vs *VillagerSuite) SetupSuite() {
	vs.playerId = types.PlayerId("1")
}

func (vs VillagerSuite) TestNewVillager() {
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
				mw.EXPECT().Poll(constants.VillagerFactionId).Return(nil)
			},
		},
		{
			name: "Ok",
			setup: func(
				mm *mock_game_logic.MockModerator,
				mw *mock_game_logic.MockWorld,
				mp *mock_game_logic.MockPoll) {
				mw.EXPECT().Poll(constants.VillagerFactionId).Return(mp).Times(2)
				mp.EXPECT().AddElectors(vs.playerId)
				mp.EXPECT().SetWeight(vs.playerId, uint(1))
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			moderator := mock_game_logic.NewMockModerator(ctrl)
			world := mock_game_logic.NewMockWorld(ctrl)
			poll := mock_game_logic.NewMockPoll(ctrl)

			moderator.EXPECT().World().Return(world)

			test.setup(moderator, world, poll)

			v, err := NewVillager(moderator, vs.playerId)

			if test.expectedErr != nil {
				vs.Nil(v)
				vs.NotNil(err)
				vs.Equal(test.expectedErr, err)
			} else {
				vs.Nil(err)
				vs.Equal(constants.VillagerRoleId, v.Id())
				vs.Equal(constants.DayPhaseId, v.(*villager).phaseId)
				vs.Equal(constants.VillagerFactionId, v.FactionId())
				vs.Equal(constants.FirstRound, v.(*villager).beginRound)
				vs.Equal(vs.playerId, v.(*villager).playerId)
				vs.Equal(constants.UnlimitedTimes, v.ActiveTimes(0))
				vs.Len(v.(*villager).abilities, 1)
				vs.Equal(action.VoteActionId, v.(*villager).abilities[0].action.Id())
				vs.True(v.(*villager).abilities[0].isImmediate)
			}
		})
	}
}

func (vs VillagerSuite) TestOnAssign() {
	ctrl := gomock.NewController(vs.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	scheduler := mock_game_logic.NewMockScheduler(ctrl)
	world := mock_game_logic.NewMockWorld(ctrl)
	poll := mock_game_logic.NewMockPoll(ctrl)

	// Mock for New fuction
	moderator.EXPECT().World().Return(world)
	moderator.EXPECT().Scheduler().Return(scheduler)
	world.EXPECT().Poll(constants.VillagerFactionId).Return(poll).Times(2)
	poll.EXPECT().AddElectors(vs.playerId)
	poll.EXPECT().SetWeight(vs.playerId, uint(1))

	v, _ := NewVillager(moderator, vs.playerId)

	moderator.EXPECT().World().Return(world).Times(3)
	world.EXPECT().Poll(constants.VillagerFactionId).Return(poll).Times(2)
	poll.EXPECT().AddElectors(vs.playerId)
	poll.EXPECT().AddCandidates(vs.playerId)
	world.EXPECT().Poll(constants.WerewolfFactionId).Return(poll)
	poll.EXPECT().AddCandidates(vs.playerId)
	scheduler.EXPECT().AddSlot(types.NewTurnSlot{
		PhaseId:    v.(*villager).phaseId,
		Turn:       v.(*villager).turn,
		BeginRound: v.(*villager).beginRound,
		PlayerId:   vs.playerId,
		RoleId:     v.Id(),
	})

	v.OnAfterAssign()
}

func (vs VillagerSuite) TestOnAfterRevoke() {
	ctrl := gomock.NewController(vs.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	scheduler := mock_game_logic.NewMockScheduler(ctrl)
	world := mock_game_logic.NewMockWorld(ctrl)
	poll := mock_game_logic.NewMockPoll(ctrl)

	// Mock for New fuction
	moderator.EXPECT().World().Return(world)
	world.EXPECT().Poll(constants.VillagerFactionId).Return(poll).Times(2)
	poll.EXPECT().AddElectors(vs.playerId)
	poll.EXPECT().SetWeight(vs.playerId, uint(1))

	v, _ := NewVillager(moderator, vs.playerId)

	moderator.EXPECT().World().Return(world).Times(3)
	moderator.EXPECT().Scheduler().Return(scheduler)
	world.EXPECT().Poll(constants.VillagerFactionId).Return(poll).Times(2)
	poll.EXPECT().RemoveElector(vs.playerId)
	poll.EXPECT().RemoveCandidate(vs.playerId)
	world.EXPECT().Poll(constants.WerewolfFactionId).Return(poll)
	poll.EXPECT().RemoveCandidate(vs.playerId)
	scheduler.EXPECT().RemoveSlot(types.RemovedTurnSlot{
		PhaseId:  v.(*villager).phaseId,
		PlayerId: vs.playerId,
		RoleId:   v.Id(),
	})

	v.OnAfterRevoke()
}
