package contract

import "uwwolf/internal/app/game/logic/types"

// Scheduler manages game's turns.
type Scheduler interface {
	// RoundID returns the latest round ID.
	RoundID() types.RoundID

	// PhaseID returns the current phase ID.
	PhaseID() types.PhaseID

	// Phase returns the current phase.
	Phase() map[types.TurnID]types.Turn

	// TurnID returns the current turn ID.
	TurnID() types.TurnID

	// Turn returns the current turn.
	Turn() types.Turn

	// CanPlay checks if the given playerID can play in the
	// current turn.
	CanPlay(playerID types.PlayerID) bool

	// PlayablePlayerIDs returns playable player ID list in
	// the current turn.
	PlayablePlayerIDs() []types.PlayerID

	// IsEmptyPhase check if specific phase is empty.
	// Check if scheduler is empty if `phaseID` is 0.
	IsEmptyPhase(phaseID types.PhaseID) bool

	// AddSlot adds new player turn to the scheduler.
	AddSlot(newSlot *types.NewTurnSlot) bool

	// RemoveSlot removes a player turn from the scheduler
	// by `TurnID` or `RoleID`.
	//
	// If `TurnID` is filled, ignore `RoleID`.
	//
	// If `PhaseID` is 0, removes all of turns of that player.
	RemoveSlot(removedSlot *types.RemovedTurnSlot) bool

	// FreezeSlot blocks slot N times.
	FreezeSlot(frozenSlot *types.FreezeTurnSlot) bool

	// NextTurn moves to the next turn.
	// Returns false if the scheduler is empty.
	NextTurn() bool
}
