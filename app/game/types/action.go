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

type KillState map[PlayerID]uint

type PredictState struct {
	Role    map[PlayerID]RoleID
	Faction map[PlayerID]FactionID
}

type RecognizeState struct {
	Role    []PlayerID
	Faction []PlayerID
}
