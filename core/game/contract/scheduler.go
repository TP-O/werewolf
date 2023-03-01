package contract

import (
	"uwwolf/game/types"
)

// Scheduler manages game's turns.
type Scheduler interface {
	// RoundID returns the latest round ID.
	RoundID() types.RoundID

	// CurrentPhaseID returns the current phase ID.
	PhaseID() types.PhaseID

	// Phase returns the current phase.
	Phase() []types.Turn

	// TurnID returns the current turn ID.
	TurnID() types.TurnID

	// Turn returns the current turn.
	Turn() types.Turn

	// PlayablePlayerIDs returns playable player ID list in
	// the current turn.
	PlayablePlayerIDs() []types.PlayerID

	// IsEmptyPhase check if specific phase is empty.
	// Check if scheduler is empty if `phaseID` is 0.
	IsEmptyPhase(phaseID types.PhaseID) bool

	// AddPlayerTurn adds new player turn to the scheduler.
	AddPlayerTurn(newTurn types.NewPlayerTurn) bool

	// RemovePlayerTurn removes a player turn from the scheduler
	// by `TurnID` or `RoleID`. Assigns -1 to `TurnID` if omitted.
	//
	// If `PhaseID` is 0, removes all of turns of that player.
	RemovePlayerTurn(removedTurn types.RemovedPlayerTurn) bool

	// NextTurn moves to the next turn.
	// Returns false if the scheduler is empty.
	NextTurn() bool
}
