package action_test

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core/action"
	"uwwolf/app/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewFactionPredict(t *testing.T) {
	predict := action.NewFactionPredict(nil, config.WerewolfFactionID)

	assert.Equal(t, config.PredictActionID, predict.ID())
	assert.NotNil(t, predict.State())
	assert.IsType(t, &types.PredictState{}, predict.State())
	assert.NotNil(t, predict.State().(*types.PredictState).Faction)
	assert.Nil(t, predict.State().(*types.PredictState).Role)

}

func TestNewRolePredict(t *testing.T) {
	predict := action.NewRolePredict(nil, config.WerewolfRoleID)

	assert.Equal(t, config.PredictActionID, predict.ID())
	assert.NotNil(t, predict.State())
	assert.IsType(t, &types.PredictState{}, predict.State())
	assert.NotNil(t, predict.State().(*types.PredictState).Role)
	assert.Nil(t, predict.State().(*types.PredictState).Faction)

}

func TestValidatePredict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)

	actorID := types.PlayerID("1")
	targetID := types.PlayerID("2")
	factionID := config.WerewolfFactionID
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr string
		setup       func(contract.Action)
	}{
		{
			name:        "Empty action request",
			req:         nil,
			expectedErr: "Action request can not be empty (⊙＿⊙')",
			setup:       func(contract.Action) {},
		},
		{
			name: "Cannot predict myself",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{actorID},
			},
			expectedErr: "WTF! You don't know who you are? (╯°□°)╯︵ ┻━┻",
			setup:       func(contract.Action) {},
		},
		{
			name: "Cannot predict known player",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
			},
			expectedErr: "You already knew this player ¯\\(º_o)/¯",
			setup: func(predict contract.Action) {
				predict.State().(*types.PredictState).Faction[targetID] = factionID
			},
		},
		{
			name: "Cannot predict non-existent player",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
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

	actorID := types.PlayerID("1")
	targetID := types.PlayerID("2")
	factionID := config.WerewolfFactionID

	mockGame.EXPECT().Player(targetID).Return(mockTarget).AnyTimes()

	tests := []struct {
		name          string
		req           *types.ActionRequest
		expectedRes   *types.ActionResponse
		expectedState map[types.PlayerID]types.FactionID
		setup         func(contract.Action)
	}{
		{
			name: "Ok with incorrect prediction",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      false,
				Message:   "",
			},
			expectedState: map[types.PlayerID]types.FactionID{
				targetID: types.FactionID(0),
			},
			setup: func(contract.Action) {
				mockTarget.EXPECT().FactionID().Return(config.VillagerFactionID)
				mockTarget.EXPECT().ID().Return(targetID)
			},
		},
		{
			name: "Ok with correct prediction",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      true,
				Message:   "",
			},
			expectedState: map[types.PlayerID]types.FactionID{
				targetID: factionID,
			},
			setup: func(contract.Action) {
				mockTarget.EXPECT().FactionID().Return(config.WerewolfFactionID)
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
			assert.Equal(t, test.expectedState, predict.State().(*types.PredictState).Faction)
		})
	}
}

func TestPerformRolePredict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)
	mockTarget := gamemock.NewMockPlayer(ctrl)

	actorID := types.PlayerID("1")
	targetID := types.PlayerID("2")
	roleID := config.WerewolfRoleID

	mockGame.EXPECT().Player(targetID).Return(mockTarget).AnyTimes()

	tests := []struct {
		name          string
		req           *types.ActionRequest
		expectedRes   *types.ActionResponse
		expectedState map[types.PlayerID]types.RoleID
		setup         func(contract.Action)
	}{
		{
			name: "Ok with incorrect prediction",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      false,
				Message:   "",
			},
			expectedState: map[types.PlayerID]types.RoleID{
				targetID: types.RoleID(0),
			},
			setup: func(contract.Action) {
				mockTarget.EXPECT().RoleIDs().Return([]types.RoleID{
					config.VillagerRoleID,
				})
				mockTarget.EXPECT().ID().Return(targetID)
			},
		},
		{
			name: "Ok with correct prediction",
			req: &types.ActionRequest{
				ActorID:   actorID,
				TargetIDs: []types.PlayerID{targetID},
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      true,
				Message:   "",
			},
			expectedState: map[types.PlayerID]types.RoleID{
				targetID: roleID,
			},
			setup: func(contract.Action) {
				mockTarget.EXPECT().RoleIDs().Return([]types.RoleID{
					config.WerewolfRoleID,
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
			assert.Equal(t, test.expectedState, predict.State().(*types.PredictState).Role)
		})
	}
}
