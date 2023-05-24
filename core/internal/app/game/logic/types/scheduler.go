package types

// Turn is type of turn index.
type Turn = uint8

// Round is type of round ID.
type Round = uint8

// TurnSlot is slot in a turn.
type TurnSlot struct {
	// BeginRound is the smallest round that the slot can be used.
	BeginRound Round

	// PlayedRound is the round that slot is used and then removed.
	// Ignore if `BeginRound` is provided.
	PlayedRound Round

	// FrozenTimes is number of remaining times the slot is blocked.
	FrozenTimes Times

	// Role is the ID of role playing in the slot.
	RoleId
}

// TurnSlots contains map players with their slot.
type TurnSlots map[PlayerId]*TurnSlot

// AddTurnSlot is added slot.
type AddTurnSlot struct {
	// PhaseId is ID of phase that the slot can be used.
	PhaseId

	// Turn is index of turn that the slot can be used.
	Turn

	// PlayerId is ID of player playing the slot.
	PlayerId

	// TurnSlot is slot data.
	TurnSlot
}

// RemoveTurnSlot is removed slot.
type RemoveTurnSlot struct {
	// PhaseId is ID of phase that the slot can be used.
	// Remove all player slots if set 0.
	PhaseId

	// Turn is index of turn that the slot can be used.
	Turn

	// PlayerId is ID of player playing the slot.
	PlayerId

	// RoleId is ID of role playing in the slot.
	// Ignored if `Turn` is provided.
	RoleId
}

// FreezeTurnSlot is frozen slot.
type FreezeTurnSlot struct {
	// PhaseId is ID of phase that the slot can be used.
	PhaseId

	// Turn is index of turn that the slot can be used.
	Turn

	// PlayerId is ID of player playing the slot.
	PlayerId

	// RoleId is ID of role playing in the slot.
	// Ignored if `TurnID` is provided.
	RoleId

	// FrozenTimes is number of remaining times the slot is blocked.
	FrozenTimes Times
}
