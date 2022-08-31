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

func TestProphecyName(t *testing.T) {
	p := action.NewProphecy(nil)

	assert.Equal(t, action.ProphecyActionName, p.Name())
}

func TestProphecyState(t *testing.T) {
	p := action.NewProphecy(nil)
	_, ok := p.State().(*state.Knowledge)

	assert.True(t, ok)
}

func TestProphecyJsonState(t *testing.T) {
	p := action.NewProphecy(nil)
	state := p.JsonState()

	assert.NotEqual(t, "{}", state)
}

func TestProphecyPerform(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)
	mockPlayer := game.NewMockPlayer(ctrl)

	//=============================================================
	p := action.NewProphecy(mockGame)
	state := p.State().(*state.Knowledge)

	werewolfId := types.PlayerId(1)
	seerId := types.PlayerId(2)

	state.Acquire(werewolfId, types.WerewolfFaction)
	state.Acquire(seerId, types.VillagerFaction)

	mockGame.
		EXPECT().
		Player(gomock.Any()).
		Return(mockPlayer).
		Times(2)

	//=============================================================
	// Actor and target is the same
	res := p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   seerId,
		Targets: []types.PlayerId{seerId},
	})

	assert.False(t, res.Ok)

	//=============================================================
	// Already known target faction
	res = p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   seerId,
		Targets: []types.PlayerId{werewolfId},
	})

	assert.False(t, res.Ok)

	//=============================================================
	// Target is werewolf
	mockPlayer.
		EXPECT().
		FactionId().
		Return(types.WerewolfFaction)

	res = p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   seerId,
		Targets: []types.PlayerId{99},
	})

	assert.True(t, res.Ok)
	assert.True(t, res.Data.(bool))

	//=============================================================
	// Target is not werewolf
	mockPlayer.
		EXPECT().
		FactionId().
		Return(types.VillagerFaction)

	res = p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   seerId,
		Targets: []types.PlayerId{100},
	})

	assert.True(t, res.Ok)
	assert.False(t, res.Data.(bool))
}
