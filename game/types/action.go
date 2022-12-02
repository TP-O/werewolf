package types

import "uwwolf/game/enum"

type ActionRequest struct {
	ActorID   enum.PlayerID
	TargetIDs []enum.PlayerID
	IsSkipped bool
}

type ActionResponse struct {
	Ok        bool
	IsSkipped bool
	Data      any
	Message   string
}

type VoteActionSetting struct {
	FactionID enum.FactionID
	PlayerID  enum.PlayerID
	Weight    uint
}

type KillState map[enum.PlayerID]uint

type PredictState struct {
	Role    map[enum.PlayerID]enum.RoleID
	Faction map[enum.PlayerID]enum.FactionID
}

type RecognizeState struct {
	Role    []enum.PlayerID
	Faction []enum.PlayerID
}
