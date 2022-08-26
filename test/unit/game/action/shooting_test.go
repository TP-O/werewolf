package action_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"uwwolf/module/game/action"
	"uwwolf/module/game/state"
	"uwwolf/test/mock/game"
	"uwwolf/types"
)

func TestShootingName(t *testing.T) {
	p := action.NewShooting(nil)

	assert.Equal(t, action.ShootingActionName, p.Name())
}

func TestShootingState(t *testing.T) {
	p := action.NewShooting(nil)
	_, ok := p.State().(*state.Shotgun)

	assert.True(t, ok)
}

func TestShootingJsonState(t *testing.T) {
	p := action.NewShooting(nil)
	state := p.JsonState()

	assert.NotEqual(t, "{}", state)
}

func TestShootingPerform(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)
	mockPlayer := game.NewMockPlayer(ctrl)

	//=============================================================
	p := action.NewShooting(mockGame)
	state := p.State().(*state.Shotgun)

	hunterId := types.PlayerId(1)
	targetId := types.PlayerId(2)

	mockGame.
		EXPECT().
		KillPlayer(gomock.Any()).
		Return(mockPlayer)

	mockPlayer.
		EXPECT().
		GetId().
		Return(targetId)

	//=============================================================
	// Actor and target is the same
	res := p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   hunterId,
		Targets: []types.PlayerId{hunterId},
	})

	assert.False(t, res.Ok)

	//=============================================================
	// Already shot
	state.Shoot(targetId)
	res = p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   hunterId,
		Targets: []types.PlayerId{targetId},
	})

	assert.False(t, res.Ok)

	//=============================================================
	// Skip shooting
	p = action.NewShooting(mockGame)

	res = p.Perform(&types.ActionRequest{
		GameId:    1,
		Actor:     hunterId,
		IsSkipped: true,
	})

	assert.True(t, res.Ok)

	//=============================================================
	// Successfully shot
	res = p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   hunterId,
		Targets: []types.PlayerId{targetId},
	})

	assert.True(t, res.Ok)
	assert.Equal(t, targetId, res.Data)
}
