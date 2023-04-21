package declare

import "uwwolf/game/types"

// One phase has 3 turn indexes by default
const (
	ZeroTurn types.TurnID = iota
	PreTurn
	MidTurn // Main turn
	PostTurn
)

// Specific ID of role turn in day phase
const (
	HunterTurnID   = PreTurn
	VillagerTurnID = MidTurn
)

// Specific ID of role turn in night phase
const (
	SeerTurnID       = PreTurn
	TwoSistersTurnID = PreTurn
	WerewolfTurnID   = MidTurn
)
