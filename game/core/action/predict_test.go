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

func TestNewFactionPredict(t *testing.T) {
	predict := action.NewFactionPredict(nil, enum.WerewolfFactionID)

	assert.Equal(t, enum.PredictActionID, predict.ID())
	assert.NotNil(t, predict.State())
	assert.IsType(t, &types.PredictState{}, predict.State())
	assert.NotNil(t, predict.State().(*types.PredictState).Faction)
	assert.Nil(t, predict.State().(*types.PredictState).Role)

}

func TestNewRolePredict(t *testing.T) {
	predict := action.NewRolePredict(nil, enum.WerewolfRoleID)

	assert.Equal(t, enum.PredictActionID, predict.ID())
	assert.NotNil(t, predict.State())
	assert.IsType(t, &types.PredictState{}, predict.State())
	assert.NotNil(t, predict.State().(*types.PredictState).Role)
	assert.Nil(t, predict.State().(*types.PredictState).Faction)

}

func TestValidatePredict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)

	actorID := enum.PlayerID("1")
	targetID := enum.PlayerID("2")
	factionID := enum.WerewolfFactionID
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr string
		setup       func(contract.Action)
	}{
		{
			name:        "Failure (Empty action request)",
			req:         nil,
			expectedErr: "Action request can not be empty (⊙＿⊙')",
			setup:       func(contract.Action) {},
		},
		{
			name: "Failure (Cannot predict myself)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{actorID},
			},
			expectedErr: "WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻",
			setup:       func(contract.Action) {},
		},
		{
			name: "Failure (Cannot predict known player)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedErr: "You already knew this player ¯\\(º_o)/¯",
			setup: func(predict contract.Action) {
				predict.State().(*types.PredictState).Faction[targetID] = factionID
			},
		},
		{
			name: "Failure (Cannot predict non-existent player)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedErr: "Non-existent player ¯\\_(ツ)_/¯",
			setup: func(contract.Action) {
				mockGame.EXPECT().Player(targetID).Return(nil).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			predict := action.NewFactionPredict(mockGame, factionID)
			test.setup(predict)
			err := predict.Validate(test.req)

			assert.Equal(t, test.expectedErr, err.Error())
		})
	}
}

func TestPerformFactionPredict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockTarget := gamemock.NewMockPlayer(ctrl)

	actorID := enum.PlayerID("1")
	targetID := enum.PlayerID("2")
	factionID := enum.WerewolfFactionID

	mockGame.EXPECT().Player(targetID).Return(mockTarget).AnyTimes()

	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
		newState    map[enum.PlayerID]enum.FactionID
		setup       func(contract.Action)
	}{
		{
			name: "Ok (Incorrect prediction)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      false,
				Message:   "",
			},
			newState: map[enum.PlayerID]enum.FactionID{
				targetID: enum.FactionID(0),
			},
			setup: func(contract.Action) {
				mockTarget.EXPECT().FactionID().Return(enum.VillagerFactionID)
				mockTarget.EXPECT().ID().Return(targetID)
			},
		},
		{
			name: "Ok (Correct prediction)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      true,
				Message:   "",
			},
			newState: map[enum.PlayerID]enum.FactionID{
				targetID: factionID,
			},
			setup: func(contract.Action) {
				mockTarget.EXPECT().FactionID().Return(enum.WerewolfFactionID)
				mockTarget.EXPECT().ID().Return(targetID)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			predict := action.NewFactionPredict(mockGame, factionID)
			test.setup(predict)
			res := predict.Perform(test.req)

			assert.Equal(t, test.expectedRes, res)
			assert.Equal(t, test.newState, predict.State().(*types.PredictState).Faction)
		})
	}
}

func TestPerformRolePredict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockTarget := gamemock.NewMockPlayer(ctrl)

	actorID := enum.PlayerID("1")
	targetID := enum.PlayerID("2")
	roleID := enum.WerewolfRoleID

	mockGame.EXPECT().Player(targetID).Return(mockTarget).AnyTimes()

	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
		newState    map[enum.PlayerID]enum.RoleID
		setup       func(contract.Action)
	}{
		{
			name: "Ok (Incorrect prediction)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      false,
				Message:   "",
			},
			newState: map[enum.PlayerID]enum.RoleID{
				targetID: enum.RoleID(0),
			},
			setup: func(contract.Action) {
				mockTarget.EXPECT().RoleIDs().Return([]enum.RoleID{
					enum.VillagerRoleID,
				})
				mockTarget.EXPECT().ID().Return(targetID)
			},
		},
		{
			name: "Ok (Correct prediction)",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []enum.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      true,
				Message:   "",
			},
			newState: map[enum.PlayerID]enum.RoleID{
				targetID: roleID,
			},
			setup: func(contract.Action) {
				mockTarget.EXPECT().RoleIDs().Return([]enum.RoleID{
					enum.WerewolfRoleID,
				})
				mockTarget.EXPECT().ID().Return(targetID)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			predict := action.NewRolePredict(mockGame, roleID)
			test.setup(predict)
			res := predict.Perform(test.req)

			assert.Equal(t, test.expectedRes, res)
			assert.Equal(t, test.newState, predict.State().(*types.PredictState).Role)
		})
	}
}
