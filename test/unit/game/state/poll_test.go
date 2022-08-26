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
	//=============================================================
	// Too few electors
	p := state.NewPoll([]types.PlayerId{1, 2})

	assert.Nil(t, p)

	//=============================================================
	// Enough electors
	p = state.NewPoll([]types.PlayerId{1, 2, 3})

	assert.NotNil(t, p)
}

func TestIsOpen(t *testing.T) {
	p := state.NewPoll(electorIds)

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

func TestPollIsAllowed(t *testing.T) {
	p := state.NewPoll(electorIds)

	//=============================================================
	// Not an elector
	assert.False(t, p.IsAllowed(999))

	//=============================================================
	// Is an elector
	assert.True(t, p.IsAllowed(electorIds[0]))

	//=============================================================
	// Is an elector who voted
	p.Open()
	p.Vote(electorIds[0], electorIds[1])

	assert.False(t, p.IsAllowed(electorIds[0]))
}

func TestGetCurrentRound(t *testing.T) {
	p := state.NewPoll(electorIds)

	//=============================================================
	// Inital
	assert.Nil(t, p.GetCurrentRound())

	//=============================================================
	// Open poll
	p.Open()

	assert.NotNil(t, p.GetCurrentRound())
}

func TestSetWeight(t *testing.T) {
	p := state.NewPoll(electorIds)

	//=============================================================
	// Valid elector id
	assert.False(t, p.SetWeight(99, 1))

	//=============================================================
	// Valid elector id
	p.SetWeight(electorIds[0], 2)
	assert.Equal(t, uint(2), p.Weights[electorIds[0]])
}

func TestOpen(t *testing.T) {
	p := state.NewPoll(electorIds)

	//=============================================================
	// Open for the first time
	assert.True(t, p.Open())
	assert.True(t, p.IsOpen())

	//=============================================================
	// Open again without closing
	assert.False(t, p.Open())
}

func TestClose(t *testing.T) {
	p := state.NewPoll(electorIds)

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
	p := state.NewPoll(electorIds)

	//=============================================================
	// Vote before opening poll
	assert.False(t, p.Vote(electorIds[0], electorIds[1]))

	//=============================================================
	// Vote with invalid elector id
	assert.False(t, p.Vote(99, electorIds[0]))

	//=============================================================
	// Successully voted
	p.Open()

	currentRound := p.GetCurrentRound()

	assert.True(t, p.Vote(electorIds[0], electorIds[1]))
	assert.NotNil(t, currentRound[electorIds[1]])
	assert.Equal(t, currentRound[electorIds[1]].Votes, uint(1))
	assert.Contains(t, currentRound[electorIds[1]].ElectorIds, electorIds[0])

	p.Close()

	//=============================================================
	// Successully voted with large weight
	p.Open()
	p.SetWeight(electorIds[0], 2)

	currentRound = p.GetCurrentRound()

	assert.True(t, p.Vote(electorIds[0], electorIds[1]))
	assert.NotNil(t, currentRound[electorIds[1]])
	assert.Equal(t, currentRound[electorIds[1]].Votes, uint(2))
	assert.Contains(t, currentRound[electorIds[1]].ElectorIds, electorIds[0])

	p.Close()

	//=============================================================
	// Vote twice
	p.Open()

	currentRound = p.GetCurrentRound()

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

func TestGetLoser(t *testing.T) {
	p := state.NewPoll(electorIds)

	//=============================================================
	// Get loser before open first poll
	assert.Equal(t, types.UnknownPlayer, p.GetLoser())

	//=============================================================
	// Get loser when poll is opening
	p.Open()

	assert.Equal(t, types.UnknownPlayer, p.GetLoser())

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

	assert.Equal(t, electorIds[1], p.GetLoser())

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

	assert.Equal(t, types.UnknownPlayer, p.GetLoser())
}

func RemoveElector(t *testing.T) {
	p := state.NewPoll(electorIds)

	//=============================================================
	// Remove non-exist elector
	assert.False(t, p.RemoveElector(99))

	//=============================================================
	// Successully removed
	assert.True(t, p.RemoveElector(electorIds[0]))
	assert.False(t, p.RemoveElector(electorIds[0]))
}
