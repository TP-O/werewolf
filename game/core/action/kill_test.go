package action_test

import (
	"testing"
	"uwwolf/game/contract"
	"uwwolf/game/core/action"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewKill(t *testing.T) {
	kill := action.NewKill(nil)

	assert.Equal(t, enum.KillActionID, kill.ID())
	assert.NotNil(t, kill.State())
	assert.IsType(t, types.KillState{}, kill.State())
}

func TestValidateKill(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockTarget := gamemock.NewMockPlayer(ctrl)

	actorID := enum.PlayerID("1")
	targetID := enum.PlayerID("2")
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
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{actorID},
			},
			expectedErr: "Appreciate your own life (｡´ ‿｀♡)",
			setup:       func() {},
		},
		{
			name: "Failure (Cannot kill non-existent player)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedErr: "Unable to kill this player!",
			setup: func() {
				mockGame.EXPECT().Player(targetID).Return(nil).Times(1)
			},
		},
		{
			name: "Failure (Cannot kill dead player)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedErr: "Unable to kill this player!",
			setup: func() {
				mockGame.EXPECT().Player(targetID).Return(mockTarget).Times(1)
				mockTarget.EXPECT().IsDead().Return(true).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			kill := action.NewKill(mockGame)
			test.setup()
			err := kill.Validate(test.req)

			assert.Equal(t, test.expectedErr, err.Error())
		})
	}
}

func TestPerformKill(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockTarget := gamemock.NewMockPlayer(ctrl)

	actorID := enum.PlayerID("1")
	targetID := enum.PlayerID("2")
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
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      targetID,
				Message:   "",
			},
			newState: types.KillState{
				targetID: 1,
			},
			setup: func(kill contract.Action) {
				mockGame.EXPECT().KillPlayer(targetID, false).Return(mockTarget).Times(1)
				mockTarget.EXPECT().ID().Return(targetID).Times(2)
			},
		},
		{
			name: "Ok (Kill second time)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      targetID,
				Message:   "",
			},
			newState: types.KillState{
				targetID: 2,
			},
			setup: func(kill contract.Action) {
				kill.State().(types.KillState)[targetID] = 1

				mockGame.EXPECT().KillPlayer(targetID, false).Return(mockTarget).Times(1)
				mockTarget.EXPECT().ID().Return(targetID).Times(2)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			kill := action.NewKill(mockGame)
			test.setup(kill)
			res := kill.Perform(test.req)

			assert.Equal(t, test.expectedRes, res)
			assert.Equal(t, test.newState[targetID], kill.State().(types.KillState)[targetID])
		})
	}
}
