package action

import (
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

// action is the basis for all concrete actions. The concrete action
// must embed this struct and modify its methods as required.
type action struct {
	// id is the action ID.
	id types.ActionID

	// game is the game instance this action affects.
	game contract.Game
}

// executable defines an interace with 3 methods required to implement the
// `Execute` method in `action“.
type executable interface {
	// validate checks if the action request is valid.
	validate(req types.ActionRequest) error

	// skip ingores the action request.
	skip(req types.ActionRequest) types.ActionResponse

	// perform completes the action request.
	perform(req types.ActionRequest) types.ActionResponse
}

// ID returns action's ID.
func (a action) ID() types.ActionID {
	return a.id
}

func (a action) skip(req types.ActionRequest) types.ActionResponse {
	return types.ActionResponse{
		Ok:      true,
		Message: "Skipped!",
	}
}

func (a action) validate(req types.ActionRequest) error {
	return nil
}

func (a action) perform(req types.ActionRequest) types.ActionResponse {
	return types.ActionResponse{
		Ok:      false,
		Message: "Nothing to do ¯\\_(ツ)_/¯",
	}
}

// execute supports the concrete action to override the `Execute` method easier
// by declaring a scheme.
func (a action) execute(exec executable, req types.ActionRequest) types.ActionResponse {
	if req.IsSkipped {
		return exec.skip(req)
	}

	if err := exec.validate(req); err != nil {
		return types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}

	return exec.perform(req)
}

// Execute checks if the request is skipped. If so, skips the action;
// otherwise, validates the request, and then performs it.
func (a action) Execute(req types.ActionRequest) types.ActionResponse {
	return a.execute(a, req)
}
