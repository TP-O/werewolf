package action

import (
	"errors"
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type PredictSuite struct {
	suite.Suite
	actorId            types.PlayerId
	targetId           types.PlayerId
	predictedFactionId types.FactionId
	predictedRoleId    types.RoleId
}

func TestPredictSuite(t *testing.T) {
	suite.Run(t, new(PredictSuite))
}

func (ps *PredictSuite) SetupSuite() {
	ps.actorId = "1"
	ps.targetId = "2"
	ps.predictedFactionId = constants.WerewolfFactionId
	ps.predictedRoleId = constants.WerewolfRoleId
}

func (ps *PredictSuite) TestNewFactionPredict() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	world := mock_game_logic.NewMockWorld(ctrl)

	pred := NewFactionPredict(world, ps.predictedFactionId).(*predict)

	ps.Equal(PredictActionId, pred.Id())
	ps.Equal(ps.predictedFactionId, pred.FactionId)
	ps.Empty(pred.Faction)
	ps.Equal(types.RoleId(0), pred.RoleId)
	ps.Empty(pred.Role)
}

func (ps *PredictSuite) TestNewRolePredict() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	world := mock_game_logic.NewMockWorld(ctrl)

	pred := NewRolePredict(world, ps.predictedRoleId).(*predict)

	ps.Equal(PredictActionId, pred.Id())
	ps.Equal(ps.predictedRoleId, pred.RoleId)
	ps.Empty(pred.Role)
	ps.Equal(types.FactionId(0), pred.FactionId)
	ps.Empty(pred.Faction)
}

func (ps *PredictSuite) TestValIdateFactionPredict() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr error
		setup       func(*predict, *mock_game_logic.MockWorld)
	}{
		{
			name: "InvalId (Cant predict yourself)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.actorId,
			},
			expectedErr: errors.New("WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻"),
			setup:       func(p *predict, mw *mock_game_logic.MockWorld) {},
		},
		{
			name: "InvalId (Cant predict known player)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.targetId,
			},
			expectedErr: errors.New("You already knew this player ¯\\(º_o)/¯"),
			setup: func(p *predict, gm *mock_game_logic.MockWorld) {
				p.Faction[ps.targetId] = true
			},
		},
		{
			name: "InvalId (Cant predict non-existent player)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: types.PlayerId("-99"),
			},
			expectedErr: errors.New("Non-existent player ¯\\_(ツ)_/¯"),
			setup: func(p *predict, gm *mock_game_logic.MockWorld) {
				gm.EXPECT().Player(types.PlayerId("-99")).Return(nil)
			},
		},
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.targetId,
			},
			setup: func(p *predict, gm *mock_game_logic.MockWorld) {
				targetPlayer := mock_game_logic.NewMockPlayer(nil)
				gm.EXPECT().Player(ps.targetId).Return(targetPlayer)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)

			pred := NewFactionPredict(world, ps.predictedFactionId).(*predict)
			test.setup(pred, world)
			err := pred.validate(test.req)

			if test.expectedErr == nil {
				ps.Nil(err)
			} else {
				ps.Equal(test.expectedErr, err)
			}
		})
	}
}

func (ps *PredictSuite) TestValIdateRolePredict() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr error
		setup       func(*predict, *mock_game_logic.MockWorld)
	}{
		{
			name: "InvalId (Cant predict yourself)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.actorId,
			},
			expectedErr: errors.New("WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻"),
			setup:       func(p *predict, mw *mock_game_logic.MockWorld) {},
		},
		{
			name: "InvalId (Cant predict known player)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.targetId,
			},
			expectedErr: errors.New("You already knew this player ¯\\(º_o)/¯"),
			setup: func(p *predict, gm *mock_game_logic.MockWorld) {
				p.Role[ps.targetId] = true
			},
		},
		{
			name: "InvalId (Cant predict non-existent player)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: types.PlayerId("-99"),
			},
			expectedErr: errors.New("Non-existent player ¯\\_(ツ)_/¯"),
			setup: func(p *predict, gm *mock_game_logic.MockWorld) {
				gm.EXPECT().Player(types.PlayerId("-99")).Return(nil)
			},
		},
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.targetId,
			},
			setup: func(p *predict, gm *mock_game_logic.MockWorld) {
				targetPlayer := mock_game_logic.NewMockPlayer(nil)
				gm.EXPECT().Player(ps.targetId).Return(targetPlayer)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)

			pred := NewRolePredict(world, ps.predictedRoleId).(*predict)
			test.setup(pred, world)
			err := pred.validate(test.req)

			if test.expectedErr == nil {
				ps.Nil(err)
			} else {
				ps.Equal(test.expectedErr, err)
			}
		})
	}
}

func (ps *PredictSuite) TestPerformFactionPredict() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes types.ActionResponse
		setup       func(*predict, *mock_game_logic.MockWorld, *mock_game_logic.MockPlayer)
	}{
		{
			name: "Ok (Incorrect prediction)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:   true,
				Data: false,
			},
			setup: func(p *predict, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mp.EXPECT().FactionId().Return(constants.VillagerFactionId)
				mp.EXPECT().Id().Return(ps.targetId).Times(2)
			},
		},
		{
			name: "Ok (Correct prediction)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:   true,
				Data: true,
			},
			setup: func(p *predict, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mp.EXPECT().FactionId().Return(ps.predictedFactionId)
				mp.EXPECT().Id().Return(ps.targetId).Times(2)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)
			targetPlayer := mock_game_logic.NewMockPlayer(ctrl)

			world.EXPECT().Player(ps.targetId).Return(targetPlayer).AnyTimes()

			pred := NewFactionPredict(world, ps.predictedFactionId).(*predict)
			test.setup(pred, world, targetPlayer)
			res := pred.perform(test.req)

			ps.Equal(test.expectedRes, res)
			ps.Equal(test.expectedRes.Data, pred.Faction[ps.targetId])
		})
	}
}

func (ps *PredictSuite) TestPerformRolePredict() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes types.ActionResponse
		setup       func(*predict, *mock_game_logic.MockWorld, *mock_game_logic.MockPlayer)
	}{
		{
			name: "Ok (Incorrect prediction)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:   true,
				Data: false,
			},
			setup: func(p *predict, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mp.EXPECT().RoleIds().Return([]types.RoleId{})
				mp.EXPECT().Id().Return(ps.targetId).Times(2)
			},
		},
		{
			name: "Ok (Correct prediction)",
			req: &types.ActionRequest{
				ActorId:  ps.actorId,
				TargetId: ps.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:   true,
				Data: true,
			},
			setup: func(p *predict, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mp.EXPECT().RoleIds().Return([]types.RoleId{ps.predictedRoleId})
				mp.EXPECT().Id().Return(ps.targetId).Times(2)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)
			targetPlayer := mock_game_logic.NewMockPlayer(ctrl)

			world.EXPECT().Player(ps.targetId).Return(targetPlayer).AnyTimes()

			pred := NewRolePredict(world, ps.predictedRoleId).(*predict)
			test.setup(pred, world, targetPlayer)
			res := pred.perform(test.req)

			ps.Equal(test.expectedRes, res)
			ps.Equal(test.expectedRes.Data, pred.Role[ps.targetId])
		})
	}
}
