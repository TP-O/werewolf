package role

import (
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type effectiveAt struct {
	IsImmediate bool

	IsRoundMatched func() bool

	IsPhaseIdMatched func() bool

	IsTurnMatched func() bool
}

// ability contains one action and its limit.
type ability struct {
	// action is a specific action.
	action contract.Action

	// activeLimit is number of times the action can be used.
	activeLimit types.Times

	effectiveAt
}
