package types

// ActionID is ID type of action
type ActionId = uint8

// ActionRequest contains the information for action execution.
type ActionRequest struct {
	// ActorId is the sender player ID.
	ActorId PlayerId

	// TargetId is the target player ID.
	TargetId PlayerId

	// IsSkipped indicates whether the request is ignored or not.
	IsSkipped bool
}

// ActionResponse contains action execution's result.
type ActionResponse struct {
	// Ok indicates whether the execution is successful or not.
	Ok bool

	// ActionId is the executed action ID.
	ActionId

	// ActionRequest is the request to execute action.
	ActionRequest

	// Message contains the error or succesful message, if any
	Message string

	// Data is the output after executing the action.
	Data any
}
