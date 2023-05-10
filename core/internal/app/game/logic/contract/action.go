package contract

import "uwwolf/internal/app/game/logic/types"

// Action executes specific action within its state.
type Action interface {
	// Id returns action's ID.
	Id() types.ActionId

	// Execute checks if the request is skipped. If so, skips the execution;
	// otherwise, validates the request, and then performs the required action.
	Execute(req types.ActionRequest) types.ActionResponse
}
