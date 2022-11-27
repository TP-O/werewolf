package contract

import "uwwolf/app/game/types"

type Action interface {
	// Id returns action id.
	ID() types.ActionID

	// State returns current action state.
	State() any

	// Perform makes some changes in state. First, it validates action request,
	// then executes it if the validation is successful. Returning struct with Ok
	// field is false means the request could not be fulfilled, otherwise execution
	// is successful.
	Perform(req *types.ActionRequest) *types.ActionResponse
}
