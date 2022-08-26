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

func TestVoteName(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)
	mockPlayer := game.NewMockPlayer(ctrl)

	mockGame.
		EXPECT().
		GetPoll(types.VillagerFaction).
		Return(&state.Poll{})

	mockPlayer.
		EXPECT().
		GetId().
		Return(types.PlayerId(1))
	mockPlayer.
		EXPECT().
		GetFactionId().
		Return(types.VillagerFaction)

	//=============================================================
	p := action.NewVote(mockGame, mockPlayer, 1)

	assert.Equal(t, action.VoteActionName, p.Name())
}

func TestVoteState(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)
	mockPlayer := game.NewMockPlayer(ctrl)

	//=============================================================
	playerId := types.PlayerId(1)

	mockGame.
		EXPECT().
		GetPoll(gomock.Any()).
		Return(state.NewPoll([]types.PlayerId{playerId, 2, 3}))

	mockPlayer.
		EXPECT().
		GetId().
		Return(playerId)
	mockPlayer.
		EXPECT().
		GetFactionId().
		Return(types.VillagerFaction)

	p := action.NewVote(mockGame, mockPlayer, 2)
	state, ok := p.State().(*state.Poll)

	//=============================================================
	// Perfect init
	assert.True(t, ok)
	assert.Equal(t, state.Weights[playerId], uint(2))
}

func TestVoteJsonState(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)
	mockPlayer := game.NewMockPlayer(ctrl)

	//=============================================================
	playerId := types.PlayerId(1)

	mockGame.
		EXPECT().
		GetPoll(gomock.Any()).
		Return(&state.Poll{})

	mockPlayer.
		EXPECT().
		GetId().
		Return(playerId)
	mockPlayer.
		EXPECT().
		GetFactionId().
		Return(types.VillagerFaction)

	p := action.NewVote(mockGame, mockPlayer, 1)
	voteState := p.JsonState()

	//=============================================================
	// Non-empty json
	assert.NotEqual(t, "{}", voteState)
}

func TestVotePerform(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)
	mockPlayer := game.NewMockPlayer(ctrl)

	//=============================================================
	playerId := types.PlayerId(1)

	mockGame.
		EXPECT().
		GetPoll(gomock.Any()).
		Return(state.NewPoll([]types.PlayerId{playerId, 2, 3})).
		AnyTimes()

	mockPlayer.
		EXPECT().
		GetId().
		Return(playerId).
		AnyTimes()
	mockPlayer.
		EXPECT().
		GetFactionId().
		Return(types.VillagerFaction).
		AnyTimes()

	p := action.NewVote(mockGame, mockPlayer, 1)
	p.State().(*state.Poll).Open()

	//=============================================================
	// Not allowed elector id
	res := p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   99,
		Targets: []types.PlayerId{playerId},
	})

	assert.False(t, res.Ok)

	//=============================================================
	// Skip vote
	res = p.Perform(&types.ActionRequest{
		GameId:    1,
		Actor:     types.PlayerId(2),
		IsSkipped: true,
	})

	assert.True(t, res.Ok)
	assert.Equal(t, types.UnknownPlayer, res.Data)

	//=============================================================
	// Successfully voted
	votedPlayer := types.PlayerId(2)

	res = p.Perform(&types.ActionRequest{
		GameId:  1,
		Actor:   playerId,
		Targets: []types.PlayerId{votedPlayer},
	})

	assert.True(t, res.Ok)
	assert.Equal(t, votedPlayer, res.Data)
}
