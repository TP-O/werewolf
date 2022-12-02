package enum

type ActionID uint

const (
	VoteActionID ActionID = iota + 1
	RecognizeActionID
	PredictActionID
	KillActionID
)
