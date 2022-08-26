package contract

import (
	"uwwolf/types"
)

// Action contains specific state and it's state will be modified
// by calling Perform.
type Action interface {
	// Name returns action name.
	Name() string

	// State returns current action state.
	State() any

	// JsonState returns current action state, but in a JSON string.
	JsonState() string

	// Perform makes some changes in state. First, it validates action request,
	// then executes it if the validation is successful. Returning struct with Ok
	// field is false means the request could not be fulfilled, otherwise execution
	// is successful.
	Perform(req *types.ActionRequest) *types.ActionResponse
}
