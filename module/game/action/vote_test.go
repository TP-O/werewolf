package action_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"uwwolf/mock/game"
	"uwwolf/module/game/action"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

func TestVoteName(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)

	mockGame.
		EXPECT().
		Poll(types.VillagerFaction).
		Return(&state.Poll{})

	//=============================================================
	playerId := types.PlayerId("1")
	factionId := types.VillagerFaction
	p := action.NewVote(mockGame, factionId, playerId, 1)

	assert.Equal(t, action.VoteActionName, p.Name())
}

func TestVoteState(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)

	//=============================================================
	playerId := types.PlayerId("1")
	poll := state.NewPoll()
	poll.AddElectors([]types.PlayerId{playerId, "2", "3"})

	mockGame.
		EXPECT().
		Poll(gomock.Any()).
		Return(poll)

	//=============================================================
	// Perfect init
	p := action.NewVote(mockGame, types.VillagerFaction, playerId, 2)
	state, ok := p.State().(*state.Poll)

	assert.True(t, ok)
	assert.Equal(t, state.Weights[playerId], uint(2))
}

func TestVotePerform(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)

	//=============================================================
	playerIds := []types.PlayerId{"1", "2"}
	poll := state.NewPoll()
	poll.AddElectors([]types.PlayerId{
		playerIds[0],
		playerIds[1],
		"3",
	})

	mockGame.
		EXPECT().
		Poll(gomock.Any()).
		Return(poll).
		Times(1)

	p := action.NewVote(mockGame, types.VillagerFaction, playerIds[0], 1)
	p.State().(*state.Poll).Open()

	//=============================================================
	// Not allowed elector id
	res := p.Perform(&types.ActionRequest{
		GameId:    1,
		ActorId:   types.PlayerId("99"),
		TargetIds: []types.PlayerId{playerIds[0]},
	})

	assert.False(t, res.Ok)

	//=============================================================
	// Skip vote
	res = p.Perform(&types.ActionRequest{
		GameId:    1,
		ActorId:   playerIds[1],
		IsSkipped: true,
	})

	assert.True(t, res.Ok)
	assert.Equal(t, types.UnknownPlayer, res.Data)

	//=============================================================
	// Successfully voted
	votedPlayer := types.PlayerId("2")

	res = p.Perform(&types.ActionRequest{
		GameId:    1,
		ActorId:   playerIds[0],
		TargetIds: []types.PlayerId{votedPlayer},
	})

	assert.True(t, res.Ok)
	assert.Equal(t, votedPlayer, res.Data)
}
