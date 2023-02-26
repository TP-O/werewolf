package contract

import (
	"uwwolf/game/types"
)

type Scheduler interface {
	// Round returns the latest round.
	RoundID() types.RoundID

	// CurrentPhaseId returns the current phase ID.
	PhaseID() types.PhaseID

	// Phase returns array of turns in the current phase.
	Phase() []*types.Turn

	// Turn returns the current turn.
	Turn() types.Turn

	// IsEmpty check if specific phase is empty.
	// Check if scheduler is empty if phaseID is 0.
	IsEmpty(phaseID types.PhaseID) bool

	// AddTurn adds new turn to the scheduler.
	AddTurn(setting types.TurnSetting) bool

	// RemoveTurn removes given role ID's the turn from the scheduler.
	// This function can make the scheduler back to previous round if
	// the removed turn is both the current turn and the current round's
	// first turn.
	RemoveTurn(roleID types.RoleID) bool

	// NextTurn moves to the next turn. If `isRemoved` is true, removes
	// the current turn before go the next one.
	// Returns false if the scheduler is empty.
	NextTurn(isRemoved bool) bool

	// FreezeTurn freezes sepecific role ID's turn in given
	// times.
	FreezeTurn(roleID types.RoleID, limit types.Limit) bool
}
