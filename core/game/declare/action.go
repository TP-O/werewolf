package declare

import "uwwolf/game/types"

// Specific action ID
const (
	VoteActionID types.ActionID = iota + 1
	IdentifyActionID
	PredictActionID
	KillActionID
)
