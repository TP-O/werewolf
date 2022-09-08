package state_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"uwwolf/module/game/state"
	"uwwolf/types"
)

var electorIds = []types.PlayerId{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
}

func TestNewPoll(t *testing.T) {
	p := state.NewPoll()

	assert.NotNil(t, p.Weights)
	assert.NotNil(t, p.Rounds)
}

func TestIsOpen(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Initial
	assert.False(t, p.IsOpen())

	//=============================================================
	// Open
	p.Open()

	assert.True(t, p.IsOpen())

	//=============================================================
	// Close
	p.Close()

	assert.False(t, p.IsOpen())
}

func TestPollIsVoted(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)
	p.Open()

	//=============================================================
	// Is an elector who didn't vote yet
	assert.False(t, p.IsVoted(electorIds[0]))

	//=============================================================
	// Is an elector who voted
	p.Vote(electorIds[0], electorIds[1])

	assert.True(t, p.IsVoted(electorIds[0]))
}

func TestPollIsAllowed(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Not an elector
	assert.False(t, p.IsAllowed(999))

	//=============================================================
	// Is an elector
	assert.True(t, p.IsAllowed(electorIds[0]))
}

func TestCurrentRound(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Inital
	assert.Nil(t, p.CurrentRound())

	//=============================================================
	// Open poll
	p.Open()

	assert.NotNil(t, p.CurrentRound())
}

func TestAddElectors(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)
	p.AddElectors([]types.PlayerId{99})

	for _, eId := range electorIds {
		assert.Contains(t, p.ElectorIds, eId)
	}

	assert.Contains(t, p.ElectorIds, types.PlayerId(99))
}

func TestSetWeight(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Valid elector id
	assert.False(t, p.SetWeight(99, 1))

	//=============================================================
	// Valid elector id
	p.SetWeight(electorIds[0], 2)
	assert.Equal(t, uint(2), p.Weights[electorIds[0]])
}

func TestOpen(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Open for the first time
	assert.True(t, p.Open())
	assert.True(t, p.IsOpen())

	//=============================================================
	// Open again without closing
	assert.False(t, p.Open())
}

func TestClose(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Close many time without opening new poll
	round1 := p.Close()
	round2 := p.Close()

	assert.Nil(t, round1)
	assert.Equal(t, round1, round2)

	//=============================================================
	// Close many time with opening new poll
	p.Open()
	round3 := p.Close()

	p.Open()
	round4 := p.Close()

	assert.False(t, p.IsOpen())
	assert.False(t, &round3 == &round4)
	assert.NotNil(t, round4[types.UnknownPlayer])
	assert.Equal(t, len(electorIds), len(round4[types.UnknownPlayer].ElectorIds))
}

func TestVote(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Vote before opening poll
	assert.False(t, p.Vote(electorIds[0], electorIds[1]))

	//=============================================================
	// Vote with invalid elector id
	assert.False(t, p.Vote(99, electorIds[0]))

	//=============================================================
	// Successully voted
	p.Open()

	currentRound := p.CurrentRound()

	assert.True(t, p.Vote(electorIds[0], electorIds[1]))
	assert.NotNil(t, currentRound[electorIds[1]])
	assert.Equal(t, currentRound[electorIds[1]].Votes, uint(1))
	assert.Contains(t, currentRound[electorIds[1]].ElectorIds, electorIds[0])

	p.Close()

	//=============================================================
	// Successully voted with large weight
	p.Open()
	p.SetWeight(electorIds[0], 2)

	currentRound = p.CurrentRound()

	assert.True(t, p.Vote(electorIds[0], electorIds[1]))
	assert.NotNil(t, currentRound[electorIds[1]])
	assert.Equal(t, currentRound[electorIds[1]].Votes, uint(2))
	assert.Contains(t, currentRound[electorIds[1]].ElectorIds, electorIds[0])

	p.Close()

	//=============================================================
	// Vote twice
	p.Open()

	currentRound = p.CurrentRound()

	p.Vote(electorIds[0], electorIds[1])

	assert.False(t, p.Vote(electorIds[0], electorIds[2]))
	assert.Nil(t, currentRound[electorIds[2]])

	p.Close()

	//=============================================================
	// Vote twice, but in difference polls
	p.Open()
	p.Vote(electorIds[0], electorIds[1])
	p.Close()
	p.Open()

	assert.True(t, p.Vote(electorIds[0], electorIds[2]))

	p.Close()
}

func TestWinner(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Get loser before open first poll
	assert.Equal(t, types.UnknownPlayer, p.Winner())

	//=============================================================
	// Get loser when poll is opening
	p.Open()

	assert.Equal(t, types.UnknownPlayer, p.Winner())

	p.Close()

	//=============================================================
	// Successully got loser - majority win
	p.Open()
	p.Vote(electorIds[0], electorIds[1])
	p.Vote(electorIds[1], electorIds[2])
	p.Vote(electorIds[2], electorIds[1])
	p.Vote(electorIds[3], electorIds[1])
	p.Vote(electorIds[4], electorIds[2])
	p.Vote(electorIds[5], electorIds[1])
	p.Vote(electorIds[6], electorIds[1])
	p.Close()

	assert.Equal(t, electorIds[1], p.Winner())

	//=============================================================
	// Successully got loser - 50/50
	p.Open()
	p.Vote(electorIds[0], electorIds[1])
	p.Vote(electorIds[1], electorIds[1])
	p.Vote(electorIds[2], electorIds[1])
	p.Vote(electorIds[3], electorIds[1])
	p.Vote(electorIds[4], electorIds[2])
	p.Vote(electorIds[5], electorIds[2])
	p.Vote(electorIds[6], electorIds[2])
	p.Vote(electorIds[7], electorIds[2])
	p.Close()

	assert.Equal(t, types.UnknownPlayer, p.Winner())
}

func RemoveElector(t *testing.T) {
	p := state.NewPoll()
	p.AddElectors(electorIds)

	//=============================================================
	// Remove non-exist elector
	assert.False(t, p.RemoveElector(99))

	//=============================================================
	// Successully removed
	assert.True(t, p.RemoveElector(electorIds[0]))
	assert.False(t, p.RemoveElector(electorIds[0]))
}
