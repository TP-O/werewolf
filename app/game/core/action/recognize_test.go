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

func TestNewFactionRecognize(t *testing.T) {
	recognize := action.NewFactionRecognize(nil, config.WerewolfFactionID)

	assert.Equal(t, config.RecognizeActionID, recognize.ID())
	assert.NotNil(t, recognize.State())
	assert.IsType(t, &types.RecognizeState{}, recognize.State())
	assert.NotNil(t, recognize.State().(*types.RecognizeState).Faction)
	assert.Nil(t, recognize.State().(*types.RecognizeState).Role)
}

func TestNewRoleRecognize(t *testing.T) {
	recognize := action.NewRoleRecognize(nil, config.WerewolfRoleID)

	assert.Equal(t, config.RecognizeActionID, recognize.ID())
	assert.NotNil(t, recognize.State())
	assert.IsType(t, &types.RecognizeState{}, recognize.State())
	assert.NotNil(t, recognize.State().(*types.RecognizeState).Role)
	assert.Nil(t, recognize.State().(*types.RecognizeState).Faction)
}

func TestPerformFactionRecognize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)

	actorID := types.PlayerID("1")
	factionID := config.WerewolfFactionID
	recognizedPlayerIDs := []types.PlayerID{"1", "2"}
	tests := []struct {
		name          string
		req           *types.ActionRequest
		expectedRes   *types.ActionResponse
		expectedState []types.PlayerID
		setup         func(contract.Action)
	}{
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorID: actorID,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      recognizedPlayerIDs,
				Message:   "",
			},
			expectedState: recognizedPlayerIDs,
			setup: func(contract.Action) {
				mockGame.EXPECT().PlayerIDsByFactionID(factionID).Return(recognizedPlayerIDs).Times(1)
			},
		},
		{
			name: "Ok but second time (return cache)",
			req: &types.ActionRequest{
				ActorID: actorID,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      recognizedPlayerIDs,
				Message:   "",
			},
			expectedState: recognizedPlayerIDs,
			setup: func(recognize contract.Action) {
				mockGame.EXPECT().PlayerIDsByFactionID(factionID).Return(recognizedPlayerIDs).Times(1)
				recognize.Execute(&types.ActionRequest{
					ActorID: actorID,
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recognize := action.NewFactionRecognize(mockGame, factionID)
			test.setup(recognize)
			res := recognize.Execute(test.req)

			assert.Equal(t, test.expectedRes, res)
			assert.Equal(t, test.expectedState, recognize.State().(*types.RecognizeState).Faction)
		})
	}
}

func TestPerformRoleRecognize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)

	actorID := types.PlayerID("1")
	roleID := config.WerewolfRoleID
	recognizedPlayerIDs := []types.PlayerID{"1", "2"}
	tests := []struct {
		name          string
		req           *types.ActionRequest
		expectedRes   *types.ActionResponse
		expectedState []types.PlayerID
		setup         func(contract.Action)
	}{
		{
			name: "Ok",
			req: &types.ActionRequest{
				ActorID: actorID,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      recognizedPlayerIDs,
				Message:   "",
			},
			expectedState: recognizedPlayerIDs,
			setup: func(contract.Action) {
				mockGame.EXPECT().PlayerIDsByRoleID(roleID).Return(recognizedPlayerIDs).Times(1)
			},
		},
		{
			name: "Ok but second time (return cache)",
			req: &types.ActionRequest{
				ActorID: actorID,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      recognizedPlayerIDs,
				Message:   "",
			},
			expectedState: recognizedPlayerIDs,
			setup: func(recognize contract.Action) {
				mockGame.EXPECT().PlayerIDsByRoleID(roleID).Return(recognizedPlayerIDs).Times(1)
				recognize.Execute(&types.ActionRequest{
					ActorID: actorID,
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recognize := action.NewRoleRecognize(mockGame, roleID)
			test.setup(recognize)
			res := recognize.Execute(test.req)

			assert.Equal(t, test.expectedRes, res)
			assert.Equal(t, test.expectedState, recognize.State().(*types.RecognizeState).Role)
		})
	}
}
