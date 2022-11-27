package config

import "uwwolf/app/game/types"

const (
	VoteActionID types.ActionID = iota + 1
	RecognizeActionID
	PredictActionID
	KillActionID
)
