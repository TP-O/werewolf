package types

// TurnID is type of turn index.
type TurnID uint8

// IsZero checks if turn ID is 0.
func (t TurnID) IsZero() bool {
	return t == 0
}

// RoundID is type of round ID.
type RoundID uint8

// IsZero checks if round ID is 0.
func (r RoundID) IsZero() bool {
	return r == 0
}

// TurnSlot is slot in a turn.
type TurnSlot struct {
	// BeginRoundID is the smallest round that the slot can be used.
	BeginRoundID RoundID

	// PlayedRoundID is the round that slot is used and then removed.
	// Ignored if `BeginRoundID` is provided.
	PlayedRoundID RoundID

	// FrozenTimes is number of remaining times the slot is blocked.
	FrozenTimes Times

	// RoleID is ID of role playing in the slot.
	RoleID
}

// Turn contains many slots.
type Turn map[PlayerID]*TurnSlot

// NewTurnSlot is added slot.
type NewTurnSlot struct {
	// PhaseID is ID of phase that the slot can be used.
	PhaseID

	// TurnID is index of turn that the slot can be used.
	TurnID

	// BeginRoundID is the smallest round that the slot can be used.
	BeginRoundID RoundID

	// PlayedRoundID is the round that slot is used and then removed.
	// Ignored if `BeginRoundID` is provided.
	PlayedRoundID RoundID

	// FrozenTimes is number of remaining times the slot is blocked.
	FrozenTimes Times

	// PlayerID is ID of player playing the slot.
	PlayerID

	// RoleID is ID of role playing in the slot.
	RoleID
}

// RemovedTurnSlot is removed slot.
type RemovedTurnSlot struct {
	// PhaseID is ID of phase that the slot can be used.
	// Remove all player slots if set 0.
	PhaseID

	// TurnID is index of turn that the slot can be used.
	TurnID

	// PlayerID is ID of player playing the slot.
	PlayerID

	// RoleID is ID of role playing in the slot.
	// Ignored if `TurnID` is provided.
	RoleID
}

// FreezeTurnSlot is frozen slot.
type FreezeTurnSlot struct {
	// PhaseID is ID of phase that the slot can be used.
	PhaseID

	// TurnID is index of turn that the slot can be used.
	TurnID

	// PlayerID is ID of player playing the slot.
	PlayerID

	// RoleID is ID of role playing in the slot.
	// Ignored if `TurnID` is provided.
	RoleID

	// FrozenTimes is number of remaining times the slot is blocked.
	FrozenTimes Times
}
