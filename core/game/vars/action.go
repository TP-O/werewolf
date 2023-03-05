package vars

import "uwwolf/game/types"

// Particular action ID
const (
	VoteActionID types.ActionID = iota + 1
	IdentifyActionID
	PredictActionID
	KillActionID
)
