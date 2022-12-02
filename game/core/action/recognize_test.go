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

func TestNewFactionRecognize(t *testing.T) {
	recognize := action.NewFactionRecognize(nil, enum.WerewolfFactionID)

	assert.Equal(t, enum.RecognizeActionID, recognize.ID())
	assert.NotNil(t, recognize.State())
	assert.IsType(t, &types.RecognizeState{}, recognize.State())
	assert.NotNil(t, recognize.State().(*types.RecognizeState).Faction)
	assert.Nil(t, recognize.State().(*types.RecognizeState).Role)
}

func TestNewRoleRecognize(t *testing.T) {
	recognize := action.NewRoleRecognize(nil, enum.WerewolfRoleID)

	assert.Equal(t, enum.RecognizeActionID, recognize.ID())
	assert.NotNil(t, recognize.State())
	assert.IsType(t, &types.RecognizeState{}, recognize.State())
	assert.NotNil(t, recognize.State().(*types.RecognizeState).Role)
	assert.Nil(t, recognize.State().(*types.RecognizeState).Faction)
}

func TestPerformFactionRecognize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)

	actorID := enum.PlayerID("1")
	factionID := enum.WerewolfFactionID
	recognizedPlayerIDs := []enum.PlayerID{"1", "2"}
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
		newState    []enum.PlayerID
		setup       func(contract.Action)
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
			newState: recognizedPlayerIDs,
			setup: func(contract.Action) {
				mockGame.EXPECT().PlayerIDsByFactionID(factionID).Return(recognizedPlayerIDs).Times(1)
			},
		},
		{
			name: "Ok (Return cache in the second time)",
			req: &types.ActionRequest{
				ActorID: actorID,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
				Data:      recognizedPlayerIDs,
				Message:   "",
			},
			newState: recognizedPlayerIDs,
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
			assert.Equal(t, test.newState, recognize.State().(*types.RecognizeState).Faction)
		})
	}
}

func TestPerformRoleRecognize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGame := gamemock.NewMockGame(ctrl)

	actorID := enum.PlayerID("1")
	roleID := enum.WerewolfRoleID
	recognizedPlayerIDs := []enum.PlayerID{"1", "2"}
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
		newState    []enum.PlayerID
		setup       func(contract.Action)
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
			newState: recognizedPlayerIDs,
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
			newState: recognizedPlayerIDs,
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
			assert.Equal(t, test.newState, recognize.State().(*types.RecognizeState).Role)
		})
	}
}
