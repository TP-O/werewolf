package action

import (
	"errors"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
)

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

func (a *action[S]) Skip(req *types.ActionRequest) *types.ActionResponse {
	return &types.ActionResponse{
		Ok:        true,
		IsSkipped: true,
	}
}

func (a *action[S]) Validate(req *types.ActionRequest) error {
	if req == nil {
		return errors.New("Action request can not be empty (⊙＿⊙')")
	}

	return nil
}

func (a *action[S]) Perform(req *types.ActionRequest) *types.ActionResponse {
	return &types.ActionResponse{
		Ok:      false,
		Message: "Nothing to do ¯\\_(ツ)_/¯",
	}
}

// combine makes `Skip`, `Validate`, and `Perform` functions work
// together.
func (a *action[S]) combine(
	skip func(*types.ActionRequest) *types.ActionResponse,
	validate func(*types.ActionRequest) error,
	perform func(*types.ActionRequest) *types.ActionResponse,
	req *types.ActionRequest,
) *types.ActionResponse {
	if req != nil && req.IsSkipped {
		return skip(req)
	}

	if err := validate(req); err != nil {
		return &types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}

	return perform(req)
}

func (a *action[S]) Execute(req *types.ActionRequest) *types.ActionResponse {
	return a.combine(a.Skip, a.Validate, a.Perform, req)
}
