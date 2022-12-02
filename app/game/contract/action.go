package contract

import "uwwolf/app/game/types"

// Action executes specific action with its state.
type Action interface {
	// ID returns action's ID.
	ID() types.ActionID

	// State returns action's current state. It records the result
	// of each successful execution, so the state can be changed after
	// each execution.
	State() any

	// Validate checks if the action request is valid.
	Validate(req *types.ActionRequest) error

	// Skip skips the task performing and may make some state changes.
	Skip(req *types.ActionRequest) *types.ActionResponse

	// Perform performs the specific task assigned to the action wihtout
	// any validation and makes some state changes if successful.
	Perform(req *types.ActionRequest) *types.ActionResponse

	// Execute combines `Skip`, `Validate`, and `Perform` in order to execute
	// the action:
	//      1. Skip
	//      2. Validate
	//      3. Perform
	Execute(req *types.ActionRequest) *types.ActionResponse
}
