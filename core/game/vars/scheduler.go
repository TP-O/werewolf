package vars

import "uwwolf/game/types"

// Round ID
const (
	ZeroRound types.RoundID = iota
	FirstRound
	SecondRound
)

// One phase has 3 turn indexes by default
const (
	ZeroTurn types.TurnID = iota
	PreTurn
	MidTurn // Main turn
	PostTurn
)
