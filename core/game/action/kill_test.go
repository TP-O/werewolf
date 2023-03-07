package action

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type KillSuite struct {
	suite.Suite
	actorID  types.PlayerID
	targetID types.PlayerID
}

func TestKillSuite(t *testing.T) {
	suite.Run(t, new(KillSuite))
}

func (ks KillSuite) SetupSuite() {
	ks.actorID = types.PlayerID("1")
	ks.targetID = types.PlayerID("2")
}

func (ks KillSuite) TestNewKill() {
	ctrl := gomock.NewController(ks.T())
	game := gamemock.NewMockGame(ctrl)

	kill := NewKill(game).(*kill)

	ks.Equal(vars.KillActionID, kill.ID())
	ks.NotNil(kill.Kills)
	ks.Len(kill.Kills, 0)
}

func (ks KillSuite) TestValidateKill() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr string
		setup       func(*gamemock.MockGame, *gamemock.MockPlayer)
	}{
		{
			name: "Invalid (Cannot commit suicide)",
			req: &types.ActionRequest{
				ActorID:  ks.actorID,
				TargetID: ks.actorID,
			},
			expectedErr: "Appreciate your own life (｡´ ‿｀♡)",
			setup:       func(mg *gamemock.MockGame, mp *gamemock.MockPlayer) {},
		},
		{
			name: "Invalid (Cannot kill non-existent player)",
			req: &types.ActionRequest{
				ActorID:  ks.actorID,
				TargetID: ks.targetID,
			},
			expectedErr: "Player does not exist (⊙＿⊙')",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPlayer) {
				mg.EXPECT().Player(ks.targetID).Return(nil).Times(1)
			},
		},
		{
			name: "Invalid (Cannot kill dead player)",
			req: &types.ActionRequest{
				ActorID:  ks.actorID,
				TargetID: ks.targetID,
			},
			expectedErr: "Player is dead [¬º-°]¬",
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPlayer) {
				mg.EXPECT().Player(ks.targetID).Return(mp).Times(1)
				mp.EXPECT().IsDead().Return(true).Times(1)
			},
		},
		{
			name: "Valid",
			req: &types.ActionRequest{
				ActorID:  ks.actorID,
				TargetID: ks.targetID,
			},
			setup: func(mg *gamemock.MockGame, mp *gamemock.MockPlayer) {
				mg.EXPECT().Player(ks.targetID).Return(mp).Times(1)
				mp.EXPECT().IsDead().Return(false).Times(1)
			},
		},
	}

	for _, test := range tests {
		ks.Run(test.name, func() {
			ctrl := gomock.NewController(ks.T())
			defer ctrl.Finish()
			game := gamemock.NewMockGame(ctrl)
			targetPlayer := gamemock.NewMockPlayer(ctrl)
			test.setup(game, targetPlayer)

			kill := NewKill(game).(*kill)
			err := kill.validate(test.req)

			if test.expectedErr == "" {
				ks.Nil(err)
			} else {
				ks.Equal(test.expectedErr, err.Error())
			}
		})
	}
}

func (ks KillSuite) TestPerformKill() {
	tests := []struct {
		name              string
		req               *types.ActionRequest
		expectedRes       *types.ActionResponse
		expectedKillTimes uint
		setup             func(*kill, *gamemock.MockGame, *gamemock.MockPlayer)
	}{
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorID:  ks.actorID,
				TargetID: ks.targetID,
			},
			expectedRes: &types.ActionResponse{
				Ok:   true,
				Data: ks.targetID,
			},
			expectedKillTimes: 1,
			setup: func(k *kill, mg *gamemock.MockGame, mp *gamemock.MockPlayer) {
				mg.EXPECT().KillPlayer(ks.targetID, false).Return(mp).Times(1)
				mp.EXPECT().ID().Return(ks.targetID).Times(2)
			},
		},
		{
			name: "Ok (Second time)",
			req: &types.ActionRequest{
				ActorID:  ks.actorID,
				TargetID: ks.targetID,
			},
			expectedRes: &types.ActionResponse{
				Ok:   true,
				Data: ks.targetID,
			},
			expectedKillTimes: 2,
			setup: func(k *kill, mg *gamemock.MockGame, mp *gamemock.MockPlayer) {
				k.Kills[ks.targetID] = 1
				mg.EXPECT().KillPlayer(ks.targetID, false).Return(mp).Times(1)
				mp.EXPECT().ID().Return(ks.targetID).Times(2)
			},
		},
	}

	for _, test := range tests {
		ks.Run(test.name, func() {
			ctrl := gomock.NewController(ks.T())
			defer ctrl.Finish()
			game := gamemock.NewMockGame(ctrl)
			targetPlayer := gamemock.NewMockPlayer(ctrl)

			kill := NewKill(game).(*kill)
			test.setup(kill, game, targetPlayer)
			res := kill.perform(test.req)

			ks.Equal(test.expectedRes, res)
			ks.Equal(test.expectedKillTimes, kill.Kills[test.req.TargetID])
		})
	}
}