package action

import (
	"encoding/json"

	"uwwolf/module/game/contract"
	"uwwolf/types"
	"uwwolf/validator"
)

type action[S any] struct {
	name  string
	state *S
	game  contract.Game
}

// Get action's name.
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
