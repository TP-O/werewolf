package action

import (
	"fmt"
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type VoteSuite struct {
	suite.Suite
	actorId   types.PlayerId
	targetId  types.PlayerId
	factionId types.FactionId
	weight    uint
}

func TestVoteSuite(t *testing.T) {
	suite.Run(t, new(VoteSuite))
}

func (vs *VoteSuite) SetupSuite() {
	vs.actorId = "1"
	vs.targetId = "2"
	vs.factionId = constants.VillagerFactionId
	vs.weight = 1
}

func (vs VoteSuite) TestNewVote() {
	tests := []struct {
		name        string
		setting     *VoteActionSetting
		expectedErr error
		setup       func(*mock_game_logic.MockWorld, *mock_game_logic.MockPoll)
	}{
		{
			name: "Failure (Poll doesnt exist)",
			setting: &VoteActionSetting{
				FactionId: vs.factionId,
			},
			expectedErr: fmt.Errorf("Poll does not exist ¯\\_(ツ)_/¯"),
			setup: func(mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPoll) {
				mw.EXPECT().Poll(vs.factionId).Return(nil)
			},
		},
		{
			name: "Ok",
			setting: &VoteActionSetting{
				FactionId: vs.factionId,
				PlayerId:  vs.actorId,
				Weight:    vs.weight,
			},
			setup: func(mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPoll) {
				mp.EXPECT().AddElectors(vs.actorId)
				mp.EXPECT().SetWeight(vs.actorId, vs.weight).Return(true)
				mw.EXPECT().Poll(vs.factionId).Return(mp).Times(2)
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)
			poll := mock_game_logic.NewMockPoll(ctrl)
			test.setup(world, poll)

			v, err := NewVote(world, test.setting)

			if test.expectedErr != nil {
				vs.Nil(v)
				vs.Equal(test.expectedErr, err)
			} else {
				vs.Equal(VoteActionId, v.Id())
				vs.NotNil(v.(*vote).poll)
				vs.Equal(poll, v.(*vote).poll)
			}
		})
	}
}

func (vs VoteSuite) TestSkip() {
	tests := []struct {
		name        string
		expectedRes types.ActionResponse
		setup       func(*mock_game_logic.MockPoll)
	}{

		{
			name: "Failure (Skip failed)",
			expectedRes: types.ActionResponse{
				Ok:      false,
				Message: "CANT_VOTE error",
			},
			setup: func(mp *mock_game_logic.MockPoll) {
				mp.EXPECT().Vote(vs.actorId, types.PlayerId("")).
					Return(false, fmt.Errorf("CANT_VOTE error"))
			},
		},
		{
			name: "Ok",
			expectedRes: types.ActionResponse{
				Ok: true,
				ActionRequest: types.ActionRequest{
					IsSkipped: true,
				},
				Message: "Skipped!",
			},
			setup: func(mp *mock_game_logic.MockPoll) {
				mp.EXPECT().Vote(vs.actorId, types.PlayerId("")).
					Return(true, nil)
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)
			poll := mock_game_logic.NewMockPoll(ctrl)

			world.EXPECT().Poll(vs.factionId).Return(poll).Times(2)
			poll.EXPECT().AddElectors(vs.actorId)
			poll.EXPECT().SetWeight(vs.actorId, vs.weight).Return(true)

			v, _ := NewVote(world, &VoteActionSetting{
				FactionId: vs.factionId,
				PlayerId:  vs.actorId,
				Weight:    vs.weight,
			})
			test.setup(poll)

			res := v.(*vote).skip(&types.ActionRequest{
				ActorId:   vs.actorId,
				IsSkipped: true,
			})

			vs.Equal(test.expectedRes, res)
		})
	}
}

func (vs VoteSuite) TestPerform() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes types.ActionResponse
		setup       func(*mock_game_logic.MockPoll)
	}{

		{
			name: "Failure (Vote failed)",
			req: &types.ActionRequest{
				ActorId:  vs.actorId,
				TargetId: vs.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:      false,
				Message: "CANT_VOTE error",
			},
			setup: func(mp *mock_game_logic.MockPoll) {
				mp.EXPECT().Vote(vs.actorId, vs.targetId).
					Return(false, fmt.Errorf("CANT_VOTE error"))
			},
		},
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorId:  vs.actorId,
				TargetId: vs.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:   true,
				Data: vs.targetId,
			},
			setup: func(mp *mock_game_logic.MockPoll) {
				mp.EXPECT().Vote(vs.actorId, vs.targetId).Return(true, nil)
			},
		},
	}

	for _, test := range tests {
		vs.Run(test.name, func() {
			ctrl := gomock.NewController(vs.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)
			poll := mock_game_logic.NewMockPoll(ctrl)

			world.EXPECT().Poll(vs.factionId).Return(poll).Times(2)
			poll.EXPECT().AddElectors(vs.actorId)
			poll.EXPECT().SetWeight(vs.actorId, vs.weight).Return(true)
			test.setup(poll)

			v, _ := NewVote(world, &VoteActionSetting{
				FactionId: vs.factionId,
				PlayerId:  vs.actorId,
				Weight:    vs.weight,
			})
			res := v.(*vote).perform(test.req)

			vs.Equal(test.expectedRes, res)
		})
	}
}
