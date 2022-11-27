package types

type ActionID uint

type ActionRequest struct {
	ActorID   PlayerID
	TargetIDs []PlayerID
	IsSkipped bool
}

type ActionResponse struct {
	Ok        bool
	IsSkipped bool
	Data      any
	Message   string
}

type VoteActionSetting struct {
	FactionID FactionID
	PlayerID  PlayerID
	Weight    uint
}
