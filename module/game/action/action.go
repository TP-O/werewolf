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

// Set action's state. Return false if type conversion is failed.
func (a *action[S]) SetState(state any) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	a.state = state.(*S)

	return true
}

// A template for embedding struct to reuse the Perform method logic.
func (a *action[S]) overridePerform(action contract.Action, data *types.ActionData) *types.PerformResult {
	if err := validator.ValidateStruct(data); err != nil {
		return &types.PerformResult{
			ErrorTag: types.InvalidInputErrorTag,
			Errors:   err,
		}
	}

	// Validate for each specific action
	if err := action.Validate(data); err != nil {
		return &types.PerformResult{
			ErrorTag: types.InvalidInputErrorTag,
			Errors:   err,
		}
	}

	return action.Execute(data)
}
