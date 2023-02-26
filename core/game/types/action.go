package types

type ActionID uint8

type ActionRequest struct {
	ActorID   PlayerID `json:"actor_id" validate:"required,len=20"`
	TargetID  PlayerID `json:"target_id" validate:"required,min=1,unique,dive,len=20"`
	IsSkipped bool     `json:"is_skipped"`
}

type StateChanges struct {
	DeadPlayerID          PlayerID
	VotedPlayerID         PlayerID
	RolePrediction        bool
	FactionPrediction     bool
	RoleIdentification    []PlayerID
	FactionIdentification []PlayerID
}

type ActionResponse struct {
	Ok           bool
	Message      string
	StateChanges StateChanges
}
