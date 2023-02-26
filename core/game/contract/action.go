package contract

import "uwwolf/game/types"

// Action executes specific action with its state.
type Action interface {
	// ID returns action's ID.
	ID() types.ActionID

	// Execute checks if the request is skipped. If so, skips the execution;
	// otherwise, validates the request, and then performs the required action.
	Execute(req types.ActionRequest) types.ActionResponse
}
