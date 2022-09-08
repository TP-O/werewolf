package action_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"uwwolf/mock/game"
	"uwwolf/module/game/action"
	"uwwolf/types"
)

func TestRecognitionName(t *testing.T) {
	p := action.NewRecognition(nil, 0)

	assert.Equal(t, action.RecognitionActionName, p.Name())
}

func TestRecognitionState(t *testing.T) {
	p := action.NewRecognition(nil, 0)

	assert.Nil(t, p.State())
}

func TestRecognitionPerform(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)

	//=============================================================
	knownRoleId := types.RoleId(1)
	sameRolePlayerIds := []types.PlayerId{1, 2, 3}

	mockGame.
		EXPECT().
		PlayerIdsWithRole(gomock.Any()).
		Return(sameRolePlayerIds).
		Times(2)

	//=============================================================
	// Already known players having same current turn's role
	r := action.NewRecognition(mockGame, knownRoleId)

	r.Perform(&types.ActionRequest{})

	res := r.Perform(&types.ActionRequest{})

	assert.False(t, res.Ok)

	//=============================================================
	// Get list of players have the same current turn's role
	r = action.NewRecognition(mockGame, knownRoleId)

	res = r.Perform(&types.ActionRequest{})

	assert.True(t, res.Ok)
	assert.Equal(t, sameRolePlayerIds, res.Data)
}
