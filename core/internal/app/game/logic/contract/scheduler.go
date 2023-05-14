package contract

import "uwwolf/internal/app/game/logic/types"

// Scheduler manages game's turns.
type Scheduler interface {
	// Round returns the current round.
	Round() types.Round

	// PhaseID returns the current phase ID.
	PhaseId() types.PhaseId

	// Phase returns the current phase.
	Phase() map[types.Turn]types.TurnSlots

	// Turn returns the current turn.
	Turn() types.Turn

	// TurnSlots returns all slots of the current turn.
	TurnSlots() types.TurnSlots

	// CanPlay checks if the given player ID can play in the
	// current turn.
	CanPlay(playerId types.PlayerId) bool

	// PlayablePlayerIds returns playable player ID list in
	// the current turn.
	PlayablePlayerIds() []types.PlayerId

	// IsEmpty check if specific phase is empty.
	// Check if scheduler is empty if `phaseId` is 0.
	IsEmpty(phaseId types.PhaseId) bool

	// AddSlot adds new player turn to the scheduler.
	AddSlot(newSlot types.AddTurnSlot) bool

	// RemoveSlot removes a player turn from the scheduler
	// by `TurnID` or `RoleID`.
	//
	// If `Turn` is provided, ignore `RoleId`.
	//
	// If `PhaseId` is 0, removes all of turns of that player.
	RemoveSlot(removeSlot types.RemoveTurnSlot) bool

	// FreezeSlot blocks slot N times.
	FreezeSlot(frozenSlot types.FreezeTurnSlot) bool

	// NextTurn moves to the next turn.
	// Returns false if the scheduler is empty.
	NextTurn() bool
}
