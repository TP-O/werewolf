package action

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

type action[S any] struct {
	name  string
	state *S
	game  contract.Game
}

type validateFnc = func(req *types.ActionRequest) string

type executeFnc = func(req *types.ActionRequest) *types.ActionResponse

func (a *action[S]) Name() string {
	return a.name
}

func (a *action[S]) State() any {
	return a.state
}

// func (a *action[S]) JsonState() string {
// 	if bytes, err := json.Marshal(a.state); err != nil {
// 		return "{}"
// 	} else {
// 		return string(bytes)
// 	}
// }

func (a *action[S]) Perform(req *types.ActionRequest) *types.ActionResponse {
	return a.perform(a.validate, a.execute, req)
}

// Execute the action request after it passes the validation.
func (a *action[S]) perform(validateFnc validateFnc, executeFnc executeFnc, req *types.ActionRequest) *types.ActionResponse {
	// Apply specific validate if general validation is passed
	alert := validateFnc(req)

	if alert != "" {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   types.InvalidInputErrorTag,
				Alert: alert,
			},
		}
	}

	return executeFnc(req)
}

// Validate the action request. Each action has different rules
// for validation. Return empty string if everything is ok.
func (a *action[S]) validate(req *types.ActionRequest) string {
	return ""
}

// Execute action request with receied data. Returning struct with Ok
// field is false means the request could not be fulfilled, otherwise execution
// is successful.
func (a *action[S]) execute(req *types.ActionRequest) *types.ActionResponse {
	return nil
}
