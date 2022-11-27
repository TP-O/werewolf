package action

import (
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
)

type validateFnc = func(req *types.ActionRequest) string

type skipFnc = func(req *types.ActionRequest) *types.ActionResponse

type executeFnc = func(req *types.ActionRequest) *types.ActionResponse

type action[S any] struct {
	state S
	id    types.ActionID
	game  contract.Game
}

func (a *action[S]) ID() types.ActionID {
	return a.id
}

func (a *action[S]) State() any {
	return a.state
}

func (a *action[S]) Perform(req *types.ActionRequest) *types.ActionResponse {
	return a.perform(a.validate, a.execute, a.skip, req)
}

// Execute the action types.ActionRequest after it passes the validation.
func (a *action[S]) perform(
	validate validateFnc,
	execute executeFnc,
	skip skipFnc,
	req *types.ActionRequest,
) *types.ActionResponse {
	if req.IsSkipped {
		return skip(req)
	}

	if msg := validate(req); msg == "" {
		return &types.ActionResponse{
			Ok:      false,
			Message: msg,
		}
	}

	return execute(req)
}

// Validate the action types.ActionRequest. Each action has different rules
// for validation. Return empty string if everything is ok.
func (a *action[S]) validate(req *types.ActionRequest) (msg string) {
	return
}

// Skip action
func (a *action[S]) skip(req *types.ActionRequest) *types.ActionResponse {
	return &types.ActionResponse{
		Ok:        true,
		IsSkipped: true,
	}
}

// Execute action types.ActionRequest with receied data. Returning struct with Ok
// field is false means the types.ActionRequest could not be fulfilled, otherwise execution
// is successful.
func (a *action[S]) execute(req *types.ActionRequest) *types.ActionResponse {
	return nil
}
