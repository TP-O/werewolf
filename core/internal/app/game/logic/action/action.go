package action

import (
	"errors"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

// executable defines an inteface with 3 required methods to implement the
// `Execute` method in `Action`.
type executable interface {
	// validate checks if the action request is valid.
	validate(req *types.ActionRequest) error

	// skip ingores the action request.
	skip(req *types.ActionRequest) types.ActionResponse

	// perform completes the action request.
	perform(req *types.ActionRequest) types.ActionResponse
}

// action is the basis for all concrete actions. The concrete action
// must embed this struct and modify its methods as required.
type action struct {
	// id is the action ID.
	id types.ActionId

	// skipValidate indicates whether validation should be skipped or not.
	skipValidate bool

	// world is the world instance affected by this action.
	world contract.World
}

var _ contract.Action = (*action)(nil)

// ID returns action's ID.
func (a action) Id() types.ActionId {
	return a.id
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (a action) Execute(req types.ActionRequest) types.ActionResponse {
	return a.execute(a, a.Id(), &req)
}

// skip ingores the action request.
func (a action) skip(req *types.ActionRequest) types.ActionResponse {
	return types.ActionResponse{
		Ok: true,
		ActionRequest: types.ActionRequest{
			IsSkipped: true,
		},
		Message: "Skipped!",
	}
}

// validate checks if the action request is valid.
func (a action) validate(req *types.ActionRequest) error {
	return errors.New("Validation is required!")
}

// perform completes the action request.
func (a action) perform(req *types.ActionRequest) types.ActionResponse {
	return types.ActionResponse{
		Ok:      false,
		Message: "Nothing to do ¯\\_(ツ)_/¯",
	}
}

// execute supports the concrete action to override the `Execute` method easier
// by declaring a schema.
func (a action) execute(exec executable, id types.ActionId, req *types.ActionRequest) types.ActionResponse {
	var res types.ActionResponse

	if req == nil {
		return types.ActionResponse{
			Ok:       false,
			ActionId: id,
			Message:  "Action request can not be empty (⊙＿⊙')",
		}
	} else if req.IsSkipped {
		res = exec.skip(req)
	} else if !a.skipValidate {
		if err := exec.validate(req); err != nil {
			res = types.ActionResponse{
				Ok:      false,
				Message: err.Error(),
			}
		} else {
			res = exec.perform(req)
		}
	} else {
		res = exec.perform(req)
	}

	res.ActionId = id
	res.ActionRequest = *req
	return res
}
