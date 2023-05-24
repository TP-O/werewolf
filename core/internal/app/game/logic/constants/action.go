package constants

import "uwwolf/internal/app/game/logic/types"

// The specific action IDs
const (
	VoteActionId types.ActionId = iota + 1
	IdentifyActionId
	PredictActionId
	KillActionId
)
