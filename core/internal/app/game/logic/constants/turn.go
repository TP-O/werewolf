package constants

import "uwwolf/internal/app/game/logic/types"

// One phase has 3 turn indexes by default
const (
	ZeroTurn types.Turn = iota
	PreTurn
	MidTurn // Main turn
	PostTurn
)

// Specific ID of role turn in day phase
const (
	HunterTurn   = PreTurn
	VillagerTurn = MidTurn
)

// Specific ID of role turn in night phase
const (
	SeerTurn       = PreTurn
	TwoSistersTurn = PreTurn
	WerewolfTurn   = MidTurn
)
