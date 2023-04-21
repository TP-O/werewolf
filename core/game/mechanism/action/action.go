package action

import (
	"uwwolf/game/mechanism/contract"
	"uwwolf/game/types"
)

// action is the basis for all concrete actions. The concrete action
// must embed this struct and modify its methods as required.
type action struct {
	// id is the action ID.
	id types.ActionID

	// world is the world instance this action affects.
	world contract.World
}

// executable defines an interace with 3 methods required to implement the
// `Execute` method in `action“.
type executable interface {
	// validate checks if the action request is valid.
	validate(req *types.ActionRequest) error

	// skip ingores the action request.
	skip(req *types.ActionRequest) *types.ActionResponse

	// perform completes the action request.
	perform(req *types.ActionRequest) *types.ActionResponse
}

// ID returns action's ID.
func (a action) ID() types.ActionID {
	return a.id
}

// skip ingores the action request.
func (a action) skip(req *types.ActionRequest) *types.ActionResponse {
	return &types.ActionResponse{
		Ok:        true,
		IsSkipped: true,
		Message:   "Skipped!",
	}
}

// validate checks if the action request is valid.
func (a action) validate(req *types.ActionRequest) error {
	return nil
}

// perform completes the action request.
func (a action) perform(req *types.ActionRequest) *types.ActionResponse {
	return &types.ActionResponse{
		Ok:      false,
		Message: "Nothing to do ¯\\_(ツ)_/¯",
	}
}

// execute supports the concrete action to override the `Execute` method easier
// by declaring a schema.
func (a action) execute(exec executable, req *types.ActionRequest) *types.ActionResponse {
	if req == nil {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Action request can not be empty (⊙＿⊙')",
		}
	} else if req.IsSkipped {
		return exec.skip(req)
	} else if err := exec.validate(req); err != nil {
		return &types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}

	return exec.perform(req)
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (a action) Execute(req *types.ActionRequest) *types.ActionResponse {
	return a.execute(a, req)
}
