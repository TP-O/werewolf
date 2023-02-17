package action

import (
	"testing"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type KillSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	game           *gamemock.MockGame
	targetedPlayer *gamemock.MockPlayer
	actorID        enum.PlayerID
	targetID       enum.PlayerID
}

func TestKillSuite(t *testing.T) {
	suite.Run(t, new(KillSuite))
}

func (ks *KillSuite) SetupSuite() {
	ks.actorID = "1"
	ks.targetID = "2"
}

func (ks *KillSuite) SetupTest() {
	ks.ctrl = gomock.NewController(ks.T())
	ks.game = gamemock.NewMockGame(ks.ctrl)
	ks.targetedPlayer = gamemock.NewMockPlayer(ks.ctrl)
}

func (ks *KillSuite) TearDownTest() {
	ks.ctrl.Finish()
}

func (ks *KillSuite) TestNewKill() {
	kill := NewKill(ks.game)

	ks.Equal(enum.KillActionID, kill.ID())
	ks.NotNil(kill.State())
	ks.IsType(types.KillState{}, kill.State())
}

func (ks *KillSuite) TestValidateKill() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr string
		setup       func()
	}{
		{
			name:        "Failure (Empty action request)",
			req:         nil,
			expectedErr: "Action request can not be empty (⊙＿⊙')",
			setup:       func() {},
		},
		{
			name: "Failure (Cannot commit suicide)",
			req: &types.ActionRequest{
				ActorID:   ks.actorID,
				TargetIDs: []enum.PlayerID{ks.actorID},
			},
			expectedErr: "Appreciate your own life (｡´ ‿｀♡)",
			setup:       func() {},
		},
		{
			name: "Failure (Cannot kill non-existent player)",
			req: &types.ActionRequest{
				ActorID:   ks.actorID,
				TargetIDs: []enum.PlayerID{ks.targetID},
			},
			expectedErr: "Player does not exist (⊙＿⊙')",
			setup: func() {
				ks.game.EXPECT().Player(ks.targetID).Return(nil).Times(1)
			},
		},
		{
			name: "Failure (Cannot kill dead player)",
			req: &types.ActionRequest{
				ActorID:   ks.actorID,
				TargetIDs: []enum.PlayerID{ks.targetID},
			},
			expectedErr: "Player is dead [¬º-°]¬",
			setup: func() {
				ks.game.EXPECT().Player(ks.targetID).Return(ks.targetedPlayer).Times(1)
				ks.targetedPlayer.EXPECT().IsDead().Return(true).Times(1)
			},
		},
	}

	for _, test := range tests {
		ks.Run(test.name, func() {
			kill := NewKill(ks.game)
			test.setup()
			err := kill.Validate(test.req)

			ks.Equal(test.expectedErr, err.Error())
		})
	}
}

func (ks *KillSuite) TestPerformKill() {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
		newState    types.KillState
		setup       func(contract.Action)
	}{
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorID:   ks.actorID,
				TargetIDs: []enum.PlayerID{ks.targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      ks.targetID,
				Message:   "",
			},
			newState: types.KillState{
				ks.targetID: 1,
			},
			setup: func(kill contract.Action) {
				ks.game.EXPECT().KillPlayer(ks.targetID, false).Return(ks.targetedPlayer).Times(1)
				ks.targetedPlayer.EXPECT().ID().Return(ks.targetID).Times(2)
			},
		},
		{
			name: "Ok (Kill second time)",
			req: &types.ActionRequest{
				ActorID:   ks.actorID,
				TargetIDs: []enum.PlayerID{ks.targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      ks.targetID,
				Message:   "",
			},
			newState: types.KillState{
				ks.targetID: 2,
			},
			setup: func(kill contract.Action) {
				kill.State().(types.KillState)[ks.targetID] = 1

				ks.game.EXPECT().KillPlayer(ks.targetID, false).Return(ks.targetedPlayer).Times(1)
				ks.targetedPlayer.EXPECT().ID().Return(ks.targetID).Times(2)
			},
		},
	}

	for _, test := range tests {
		ks.Run(test.name, func() {
			kill := NewKill(ks.game)
			test.setup(kill)
			res := kill.Perform(test.req)

			ks.Equal(test.expectedRes, res)
			ks.Equal(test.newState[ks.targetID], kill.State().(types.KillState)[ks.targetID])
		})
	}
}
