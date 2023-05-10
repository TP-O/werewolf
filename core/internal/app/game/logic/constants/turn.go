package constants

import "uwwolf/internal/app/game/logic/types"

// One phase has 3 turn indexes by default
const (
	ZeroTurn types.TurnId = iota
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
