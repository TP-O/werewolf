package types

// TurnID is type of turn index.
type Turn = uint8

// RoundID is type of round ID.
type Round = uint8

// TurnSlot is slot in a turn.
type TurnSlot struct {
	// BeginRoundID is the smallest round that the slot can be used.
	BeginRound Round

	// PlayedRoundID is the round that slot is used and then removed.
	// Ignored if `BeginRoundID` is provided.
	PlayedRound Round

	// FrozenTimes is number of remaining times the slot is blocked.
	FrozenTimes Times

	// RoleID is ID of role playing in the slot.
	RoleId
}

// Turn contains many slots.
type TurnRecords map[PlayerId]*TurnSlot

// NewTurnSlot is added slot.
type NewTurnSlot struct {
	// PhaseID is ID of phase that the slot can be used.
	PhaseId

	// TurnID is index of turn that the slot can be used.
	Turn

	// BeginRoundID is the smallest round that the slot can be used.
	BeginRound Round

	// PlayedRoundID is the round that slot is used and then removed.
	// Ignored if `BeginRoundID` is provided.
	PlayedRound Round

	// FrozenTimes is number of remaining times the slot is blocked.
	FrozenTimes Times

	// PlayerID is ID of player playing the slot.
	PlayerId

	// RoleID is ID of role playing in the slot.
	RoleId
}

// RemovedTurnSlot is removed slot.
type RemovedTurnSlot struct {
	// PhaseID is ID of phase that the slot can be used.
	// Remove all player slots if set 0.
	PhaseId

	// TurnID is index of turn that the slot can be used.
	Turn

	// PlayerID is ID of player playing the slot.
	PlayerId

	// RoleID is ID of role playing in the slot.
	// Ignored if `TurnID` is provided.
	RoleId
}

// FreezeTurnSlot is frozen slot.
type FreezeTurnSlot struct {
	// PhaseID is ID of phase that the slot can be used.
	PhaseId

	// TurnID is index of turn that the slot can be used.
	Turn

	// PlayerID is ID of player playing the slot.
	PlayerId

	// RoleID is ID of role playing in the slot.
	// Ignored if `TurnID` is provided.
	RoleId

	// FrozenTimes is number of remaining times the slot is blocked.
	FrozenTimes Times
}
