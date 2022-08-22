package state_test

import (
	"testing"

	"uwwolf/module/game/state"
	"uwwolf/types"

	"github.com/stretchr/testify/assert"
)

var turnSettings []*types.TurnSetting = []*types.TurnSetting{
	{
		PhaseId:    types.NightPhase,
		RoleId:     1,
		PlayerIds:  []types.PlayerId{1},
		BeginRound: 1,
		Priority:   1,
		Expiration: types.UnlimitedTimes,
	},
	{
		PhaseId:    types.DayPhase,
		RoleId:     2,
		PlayerIds:  []types.PlayerId{2},
		BeginRound: 1,
		Priority:   1,
		Expiration: types.UnlimitedTimes,
	},
	{
		PhaseId:    types.DuskPhase,
		RoleId:     3,
		PlayerIds:  []types.PlayerId{3},
		BeginRound: 1,
		Priority:   1,
		Expiration: types.UnlimitedTimes,
	},
}

func TestGetCurrentTurn(t *testing.T) {
	r := state.NewRound()

	//=============================================================
	// Empty round
	assert.Nil(t, r.GetCurrentTurn())

	//=============================================================
	// Non-empty round
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.Equal(t, turnSettings[0].RoleId, r.GetCurrentTurn().RoleId())

	r.NextTurn()

	assert.Equal(t, turnSettings[1].RoleId, r.GetCurrentTurn().RoleId())

	r.NextTurn()

	assert.Equal(t, turnSettings[2].RoleId, r.GetCurrentTurn().RoleId())
}

func TestReset(t *testing.T) {
	r := state.NewRound()

	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	r.Reset()

	assert.Equal(t, types.RoundId(1), r.GetCurrentId())
	assert.Equal(t, 0, len(r.GetCurrentPhase()))
	assert.Nil(t, r.GetCurrentTurn())
	assert.True(t, r.IsEmpty())
}

func TestRoundIsAllowed(t *testing.T) {
	r := state.NewRound()
	playerId := types.PlayerId(1)

	//=============================================================
	// Empty round
	assert.False(t, r.IsAllowed(playerId))

	r.Reset()

	//=============================================================
	// Invalid player id
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	r.NextTurn()

	assert.False(t, r.IsAllowed(playerId))

	r.Reset()

	//=============================================================
	// Valid player id
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.True(t, r.IsAllowed(playerId))

	r.Reset()
}

func TestIsEmpty(t *testing.T) {
	r := state.NewRound()

	//=============================================================
	// Empty round
	assert.True(t, r.IsEmpty())

	r.Reset()

	//=============================================================
	// Non-empty round
	r.AddTurn(turnSettings[0])

	assert.False(t, r.IsEmpty())

	r.Reset()
}

func TestNextTurn(t *testing.T) {
	r := state.NewRound()

	//=============================================================
	// Empty round
	assert.False(t, r.NextTurn())

	r.Reset()

	//=============================================================
	// Next turn with changing phase
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	firstPhase := r.GetCurrentPhase()
	r.NextTurn()
	secondPhase := r.GetCurrentPhase()

	assert.NotEqual(t, firstPhase, secondPhase)

	r.Reset()

	//=============================================================
	// Next turn without changing phase
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	r.AddTurn(&types.TurnSetting{
		PhaseId:  types.NightPhase,
		Position: types.LastPosition,
	})

	unchagedPhase := r.GetCurrentPhase()
	r.NextTurn()

	assert.Equal(t, unchagedPhase, r.GetCurrentPhase())

	r.Reset()

	//=============================================================
	// Delete turn if it's expired
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	oneTimesTurnSetting := &types.TurnSetting{
		RoleId:     99,
		PhaseId:    types.NightPhase,
		Expiration: types.OneTimes,
		Position:   types.LastPosition,
	}

	r.AddTurn(oneTimesTurnSetting)

	r.NextTurn()

	assert.Equal(t, oneTimesTurnSetting.RoleId, r.GetCurrentTurn().RoleId())

	r.NextTurn()
	r.NextTurn()
	r.NextTurn()
	r.NextTurn()

	assert.Equal(t, turnSettings[1].RoleId, r.GetCurrentTurn().RoleId())
	assert.NotEqual(t, oneTimesTurnSetting.RoleId, r.GetCurrentTurn().RoleId())

	r.Reset()
}

func TestAddTurn(t *testing.T) {
	r := state.NewRound()

	//=============================================================
	// Add turn with invalid phase id
	assert.False(t, r.AddTurn(&types.TurnSetting{
		PhaseId: 99,
	}))
	assert.True(t, r.IsEmpty())

	r.Reset()

	//=============================================================
	// Add turn with invalid position id
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.False(t, r.AddTurn(&types.TurnSetting{
		PhaseId:  types.NightPhase,
		Position: -5,
	}))
	assert.False(t, r.AddTurn(&types.TurnSetting{
		PhaseId:  types.NightPhase,
		Position: 99,
	}))

	r.Reset()

	//=============================================================
	// Add turn with existed role id
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.False(t, r.AddTurn(&types.TurnSetting{
		RoleId:  turnSettings[0].RoleId,
		PhaseId: types.NightPhase,
	}))

	r.Reset()

	//=============================================================
	// Add turn to empty phases
	for _, setting := range turnSettings {
		assert.True(t, r.AddTurn(setting))
	}

	firstPhase := r.GetCurrentPhase()

	assert.Equal(t, 1, len(firstPhase))
	assert.False(t, r.IsEmpty())

	r.NextTurn()
	secondPhase := r.GetCurrentPhase()

	assert.Equal(t, 1, len(secondPhase))
	assert.NotEqual(t, firstPhase, secondPhase)

	r.NextTurn()
	thirdPhase := r.GetCurrentPhase()

	assert.Equal(t, 1, len(thirdPhase))
	assert.NotEqual(t, firstPhase, thirdPhase)
	assert.NotEqual(t, secondPhase, thirdPhase)

	r.Reset()

	//=============================================================
	// Add turn to next position
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	nextTurn := &types.TurnSetting{
		PhaseId:  types.NightPhase,
		RoleId:   99,
		Position: types.NextPosition,
	}

	r.AddTurn(nextTurn)
	r.NextTurn()

	assert.Equal(t, nextTurn.RoleId, r.GetCurrentTurn().RoleId())

	r.Reset()

	//=============================================================
	// Add turn to sorted position
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	sortedTurn1 := &types.TurnSetting{
		PhaseId:  types.NightPhase,
		RoleId:   99,
		Priority: 3,
		Position: types.SortedPosition,
	}
	sortedTurn2 := &types.TurnSetting{
		PhaseId:  types.NightPhase,
		RoleId:   100,
		Priority: 4,
		Position: types.SortedPosition,
	}
	sortedTurn3 := &types.TurnSetting{
		PhaseId:  types.NightPhase,
		RoleId:   101,
		Priority: 5,
		Position: types.SortedPosition,
	}

	r.AddTurn(sortedTurn3)
	r.AddTurn(sortedTurn2)
	r.AddTurn(sortedTurn1)

	r.NextTurn()

	assert.Equal(t, sortedTurn1.RoleId, r.GetCurrentTurn().RoleId())

	r.NextTurn()

	assert.Equal(t, sortedTurn2.RoleId, r.GetCurrentTurn().RoleId())

	r.NextTurn()

	assert.Equal(t, sortedTurn3.RoleId, r.GetCurrentTurn().RoleId())

	r.Reset()

	//=============================================================
	// Add turn to last position
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	r.AddTurn(&types.TurnSetting{
		PhaseId:  types.NightPhase,
		RoleId:   99,
		Position: types.LastPosition,
	})

	lastTurn := &types.TurnSetting{
		PhaseId:  types.NightPhase,
		RoleId:   100,
		Position: types.LastPosition,
	}

	r.AddTurn(lastTurn)
	currentTurn01 := r.GetCurrentTurn()

	assert.Equal(t, currentTurn01, r.GetCurrentTurn())

	r.NextTurn()
	r.NextTurn()

	assert.Equal(t, lastTurn.RoleId, r.GetCurrentTurn().RoleId())

	r.Reset()

	//=============================================================
	// Add turn to other positions
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	r.AddTurn(&types.TurnSetting{
		PhaseId:  types.NightPhase,
		RoleId:   99,
		Position: 1,
	})

	otherTurn := &types.TurnSetting{
		PhaseId:  types.NightPhase,
		RoleId:   100,
		Position: 1,
	}

	r.AddTurn(otherTurn)

	r.NextTurn()

	assert.Equal(t, otherTurn.RoleId, r.GetCurrentTurn().RoleId())

	r.Reset()

	//=============================================================
	// Add new turn has index less than or equal current turn
	// index, but current turn is unchanged
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	// Equal to current index
	currentTurn02 := r.GetCurrentTurn()

	r.AddTurn(&types.TurnSetting{
		PhaseId:  types.NightPhase,
		Position: 0,
	})

	assert.Equal(t, currentTurn02, r.GetCurrentTurn())

	// Less than current index
	r.NextTurn()
	currentTurn03 := r.GetCurrentTurn()

	r.AddTurn(&types.TurnSetting{
		PhaseId:  types.NightPhase,
		Position: 0,
	})

	assert.Equal(t, currentTurn03, r.GetCurrentTurn())

	r.Reset()
}

func TestRemoveTurn(t *testing.T) {
	r := state.NewRound()

	//=============================================================
	// Non-exist role id
	assert.False(t, r.RemoveTurn(99))

	r.Reset()

	//=============================================================
	// Remove turn in normal case
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	currentPhase01 := r.GetCurrentPhase()
	removedTurnSetting := &types.TurnSetting{
		RoleId:   99,
		PhaseId:  types.NightPhase,
		Position: types.LastPosition,
	}

	r.AddTurn(removedTurnSetting)
	r.NextTurn()

	assert.True(t, r.RemoveTurn(removedTurnSetting.RoleId))
	assert.Equal(t, currentPhase01, r.GetCurrentPhase())

	r.Reset()

	//=============================================================
	// Remove current turn
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	currentTurn01 := r.GetCurrentTurn()
	r.RemoveTurn(turnSettings[0].RoleId)

	assert.NotEqual(t, currentTurn01, r.GetCurrentTurn())
	assert.Equal(t, turnSettings[2].RoleId, r.GetCurrentTurn().RoleId())

	r.Reset()

	//=============================================================
	// Remove turn has index less than current turn index
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	r.AddTurn(&types.TurnSetting{
		PhaseId:  types.NightPhase,
		Position: types.LastPosition,
	})
	r.NextTurn()

	currentTurn02 := r.GetCurrentTurn()

	r.RemoveTurn(turnSettings[0].RoleId)

	assert.Equal(t, currentTurn02, r.GetCurrentTurn())

	r.Reset()
}

func TestAddPlayer(t *testing.T) {
	r := state.NewRound()

	//=============================================================
	// Add to non-exist role
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.False(t, r.AddPlayer(1, 99))

	r.Reset()

	//=============================================================
	// Add existed player id
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.False(t, r.AddPlayer(turnSettings[0].PlayerIds[0], turnSettings[0].RoleId))
	assert.Len(t, r.GetCurrentTurn().PlayerIds(), 1)

	r.Reset()

	//=============================================================
	// Add successully
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.True(t, r.AddPlayer(types.PlayerId(2), turnSettings[0].RoleId))
	assert.Contains(t, r.GetCurrentTurn().PlayerIds(), types.PlayerId(2))

	r.Reset()
}

func TestDeletePlayer(t *testing.T) {
	r := state.NewRound()

	//=============================================================
	// Delete to non-exist role
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.False(t, r.DeletePlayer(1, 99))

	r.Reset()

	//=============================================================
	// Delete successully
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	assert.True(t, r.DeletePlayer(types.PlayerId(1), turnSettings[0].RoleId))
	assert.NotContains(t, r.GetCurrentTurn().PlayerIds(), types.PlayerId(1))

	r.Reset()
}

func TestDeletePlayerFromAllTurns(t *testing.T) {
	r := state.NewRound()

	//=============================================================
	// Delete successfully
	for _, setting := range turnSettings {
		r.AddTurn(setting)
	}

	r.AddTurn(&types.TurnSetting{
		PhaseId:   types.DayPhase,
		PlayerIds: []types.PlayerId{1},
		Position:  0,
	})

	r.DeletePlayerFromAllTurns(types.PlayerId(1))

	assert.NotContains(t, r.GetCurrentTurn().PlayerIds(), types.PlayerId(1))

	r.NextTurn()

	assert.NotContains(t, r.GetCurrentTurn().PlayerIds(), types.PlayerId(1))

	r.Reset()
}
