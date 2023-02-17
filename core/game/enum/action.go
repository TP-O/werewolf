package enum

type ActionID = uint32

const (
	VoteActionID ActionID = iota + 1
	RecognizeActionID
	PredictActionID
	KillActionID
)
