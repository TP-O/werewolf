package role

import (
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type effectiveAt struct {
	isImmediate bool

	isRoundMatched func() bool

	isPhaseIdMatched func() bool

	isTurnMatched func() bool
}

func (e effectiveAt) CanExecute() bool {
	return e.isRoundMatched() && e.isPhaseIdMatched() && e.isTurnMatched()
}

// ability contains one action and its limit.
type ability struct {
	// action is a specific action.
	action contract.Action

	// activeLimit is number of times the action can be used.
	activeLimit types.Times

	effectiveAt
}
