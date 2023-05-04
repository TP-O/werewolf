package types

// ActionID is ID type of action
type ActionID uint8

// ActionRequest contains information for action execution.
type ActionRequest struct {
	// ActorID is player ID of request sender.
	ActorID PlayerID

	// TargetID is player ID of target player.
	TargetID PlayerID

	// IsSkipped marks that the request is ignored.
	IsSkipped bool
}

// ActionResponse contains action execution's result.
type ActionResponse struct {
	// Ok marks that is execution is successful.
	Ok bool

	// RoundID is round ID which the action is executed.
	RoundID

	// RoleID is ID of role executing the action.
	RoleID

	// ActionID is executed action ID.
	ActionID

	// TargetID is player ID of affected player.
	TargetID PlayerID

	// IsSkipped marks that the request is ignored.
	IsSkipped bool

	// Message contains error or succesful message, if any
	Message string

	// Data is expected data when executing the action.
	Data any
}
