package contract

import "uwwolf/app/game/types"

type Scheduler interface {
	// CurrentId returns the current round id.
	Round() types.Round

	// CurrentPhaseId returns the current phase id.
	PhaseID() types.PhaseID

	Phase() []*types.Turn

	// CurrentTurn returns the current turn instance.
	Turn() *types.Turn

	// IsEmpty returns true if all phases in round are empty.
	IsEmpty(phaseID types.PhaseID) bool

	// AddTurn adds new turn to the round. Returns false if
	// new turn is invalid.
	AddTurn(setting *types.TurnSetting) bool

	// RemoveTurn removes the turn of given role id from
	// the round.
	RemoveTurn(roleID types.RoleID) bool

	// NextTurn moves to the next turn. Returns false if
	// the round is empty.
	NextTurn(isRemoved bool) bool

	FreezeTurn(roleID types.RoleID, limit types.Limit) bool
}
