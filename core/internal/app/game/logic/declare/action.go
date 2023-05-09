package declare

import "uwwolf/internal/app/game/logic/types"

// Specific action ID
const (
	VoteActionID types.ActionID = iota + 1
	IdentifyActionID
	PredictActionID
	KillActionID
)
