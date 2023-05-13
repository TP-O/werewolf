package contract

import "uwwolf/internal/app/game/logic/types"

// Scheduler manages game's turns.
type Scheduler interface {
	// RoundID returns the latest round ID.
	Round() types.Round

	// PhaseID returns the current phase ID.
	PhaseId() types.PhaseId

	// Phase returns the current phase.
	Phase() map[types.Turn]types.TurnRecords

	// TurnID returns the current turn ID.
	Turn() types.Turn

	// Turn returns the current turn.
	TurnRecords() types.TurnRecords

	// CanPlay checks if the given playerID can play in the
	// current turn.
	CanPlay(playerID types.PlayerId) bool

	// PlayablePlayerIDs returns playable player ID list in
	// the current turn.
	PlayablePlayerIDs() []types.PlayerId

	// IsEmptyPhase check if specific phase is empty.
	// Check if scheduler is empty if `phaseID` is 0.
	IsEmptyPhase(phaseId types.PhaseId) bool

	// AddSlot adds new player turn to the scheduler.
	AddSlot(newSlot types.NewTurnSlot) bool

	// RemoveSlot removes a player turn from the scheduler
	// by `TurnID` or `RoleID`.
	//
	// If `TurnID` is filled, ignore `RoleID`.
	//
	// If `PhaseID` is 0, removes all of turns of that player.
	RemoveSlot(removedSlot types.RemovedTurnSlot) bool

	// FreezeSlot blocks slot N times.
	FreezeSlot(frozenSlot types.FreezeTurnSlot) bool

	// NextTurn moves to the next turn.
	// Returns false if the scheduler is empty.
	NextTurn() bool
}
