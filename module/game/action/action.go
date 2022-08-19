package action

import (
	"encoding/json"

	govalidator "github.com/go-playground/validator/v10"

	"uwwolf/module/game/contract"
	"uwwolf/types"
	"uwwolf/validator"
)

type action[S any] struct {
	name  string
	state *S
	game  contract.Game
}

func (a *action[S]) GetName() string {
	return a.name
}

// Export action's state  as JSON string.
func (a *action[S]) JsonState() string {
	if bytes, err := json.Marshal(a.state); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}

// A template for embedding struct to reuse the Perform method logic.
func (a *action[S]) overridePerform(action contract.Action, req *types.ActionRequest) *types.ActionResponse {
	err := validator.ValidateStruct(req)

	// Apply specific validate if general validation is passed
	if err == nil {
		err = action.Validate(req)
	}

	if err != nil {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag: types.InvalidInputErrorTag,
				Msg: err,
			},
		}
	}

	return action.Execute(req)
}

// Validate action's input first, then execute it if the
// validation is successful. Only supposed to fail if
// and only if an error message is returned.
func (a *action[S]) Perform(req *types.ActionRequest) *types.ActionResponse {
	return &types.ActionResponse{}
}

// Validate the action's input. Each action has different rules
// for data validation.
func (a *action[S]) Validate(req *types.ActionRequest) govalidator.ValidationErrorsTranslations {
	return nil
}

// Execute the action with receied data. Return the result of execution
// and error message, if any. The execution is only supposed to fail if
// and only if an error message is returned. The first response arg is
// just a status of the execution, so its meaning depends on contenxt.
func (a *action[S]) Execute(req *types.ActionRequest) *types.ActionResponse {
	return &types.ActionResponse{}
}
