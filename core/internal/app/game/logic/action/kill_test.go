package action

import (
	"fmt"
	"testing"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type KillSuite struct {
	suite.Suite
	actorId  types.PlayerId
	targetId types.PlayerId
}

func TestKillSuite(t *testing.T) {
	suite.Run(t, new(KillSuite))
}

func (ks *KillSuite) SetupSuite() {
	ks.actorId = "1"
	ks.targetId = "2"
}

func (ks KillSuite) TestNewKill() {
	ctrl := gomock.NewController(ks.T())
	world := mock_game_logic.NewMockWorld(ctrl)

	kill := NewKill(world).(*kill)

	ks.Equal(KillActionId, kill.Id())
	ks.NotNil(kill.Kills)
	ks.Empty(kill.Kills)
}

func (ks KillSuite) TestValidate() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr error
		setup       func(*mock_game_logic.MockWorld, *mock_game_logic.MockPlayer)
	}{
		{
			name: "Invalid (Cannot commit suicide)",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.actorId,
			},
			expectedErr: fmt.Errorf("Appreciate your own life (｡´ ‿｀♡)"),
			setup:       func(mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {},
		},
		{
			name: "Invalid (Cannot kill non-existent player)",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.targetId,
			},
			expectedErr: fmt.Errorf("Player does not exist (⊙＿⊙')"),
			setup: func(mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mw.EXPECT().Player(ks.targetId).Return(nil)
			},
		},
		{
			name: "Invalid (Cannot kill dead player)",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.targetId,
			},
			expectedErr: fmt.Errorf("Player is dead [¬º-°]¬"),
			setup: func(mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mw.EXPECT().Player(ks.targetId).Return(mp)
				mp.EXPECT().IsDead().Return(true)
			},
		},
		{
			name: "Valid",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.targetId,
			},
			setup: func(mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mw.EXPECT().Player(ks.targetId).Return(mp)
				mp.EXPECT().IsDead().Return(false)
			},
		},
	}

	for _, test := range tests {
		ks.Run(test.name, func() {
			ctrl := gomock.NewController(ks.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)
			targetPlayer := mock_game_logic.NewMockPlayer(ctrl)
			test.setup(world, targetPlayer)

			kill := NewKill(world).(*kill)
			err := kill.validate(test.req)

			if test.expectedErr == nil {
				ks.Nil(err)
			} else {
				ks.Equal(test.expectedErr, err)
			}
		})
	}
}

func (ks KillSuite) TestPerform() {
	tests := []struct {
		name              string
		req               *types.ActionRequest
		expectedRes       types.ActionResponse
		expectedKillTimes uint
		setup             func(*kill, *mock_game_logic.MockWorld, *mock_game_logic.MockPlayer)
	}{
		{
			name: "Failure (Target doesn't exist)",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:      false,
				Message: "The targeted player doesn't exist!",
			},
			expectedKillTimes: 0,
			setup: func(k *kill, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mw.EXPECT().Player(ks.targetId).Return(nil)
			},
		},
		{
			name: "Ok (Kill covered player)",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:      false,
				Message: "Unable to kill the targeted player!",
			},
			expectedKillTimes: 0,
			setup: func(k *kill, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mw.EXPECT().Player(ks.targetId).Return(mp)
				mp.EXPECT().Die().Return(false)
			},
		},
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:   true,
				Data: ks.targetId,
			},
			expectedKillTimes: 1,
			setup: func(k *kill, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mw.EXPECT().Player(ks.targetId).Return(mp)
				mp.EXPECT().Die().Return(true)
				mp.EXPECT().Id().Return(ks.targetId).Times(2)
			},
		},
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:   true,
				Data: ks.targetId,
			},
			expectedKillTimes: 1,
			setup: func(k *kill, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				mw.EXPECT().Player(ks.targetId).Return(mp)
				mp.EXPECT().Die().Return(true)
				mp.EXPECT().Id().Return(ks.targetId).Times(2)
			},
		},
		{
			name: "Ok (Second time)",
			req: &types.ActionRequest{
				ActorId:  ks.actorId,
				TargetId: ks.targetId,
			},
			expectedRes: types.ActionResponse{
				Ok:   true,
				Data: ks.targetId,
			},
			expectedKillTimes: 2,
			setup: func(k *kill, mw *mock_game_logic.MockWorld, mp *mock_game_logic.MockPlayer) {
				k.Kills[ks.targetId] = 1
				mw.EXPECT().Player(ks.targetId).Return(mp)
				mp.EXPECT().Die().Return(true)
				mp.EXPECT().Id().Return(ks.targetId).Times(2)
			},
		},
	}

	for _, test := range tests {
		ks.Run(test.name, func() {
			ctrl := gomock.NewController(ks.T())
			defer ctrl.Finish()
			world := mock_game_logic.NewMockWorld(ctrl)
			targetPlayer := mock_game_logic.NewMockPlayer(ctrl)

			kill := NewKill(world).(*kill)
			test.setup(kill, world, targetPlayer)
			res := kill.perform(test.req)

			ks.Equal(test.expectedRes, res)
			ks.Equal(test.expectedKillTimes, kill.Kills[test.req.TargetId])
		})
	}
}
