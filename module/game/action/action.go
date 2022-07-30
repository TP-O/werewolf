package action

import (
	"encoding/json"
	"errors"

	"uwwolf/module/game/core"
	"uwwolf/types"
	"uwwolf/validator"
)

type action[S any] struct {
	name  string
	state *S
	game  core.Game
}

type Action interface {
	// Get action's name.
	GetName() string

	// Export action's state  as JSON string.
	JsonState() string

	// Set action's state. Return false if type conversion is failed.
	SetState(state any) bool

	// Validate action's input first, then execute it if the
	// validation is successful. Only supposed to fail if
	// and only if an error message is returned.
	Perform(data *types.ActionData) (bool, error)

	// Validate the action's input. Each action has different rules
	// for data validation.
	validate(data *types.ActionData) (bool, error)

	// Execute the action with receied data. Return the result of execution
	// and error message, if any. The execution is only supposed to fail if
	// and only if an error message is returned. The first response arg is
	// just a status of the execution, so its meaning depends on contenxt.
	execute(data *types.ActionData) (bool, error)
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
func (a *action[S]) overridePerform(action Action, data *types.ActionData) (bool, error) {
	if !validator.SimpleValidateStruct(data) {
		return false, errors.New("Invalid action!")
	}

	// Validate for each specific action
	if _, err := action.validate(data); err != nil {
		return false, err
	}

	return action.execute(data)
}

// Validate action's input first, then execute it if the
// validation is successful. Only supposed to fail if
// and only if an error message is returned.
func (a *action[S]) Perform(data *types.ActionData) (bool, error)

// Validate the action's input. Each action has different rules
// for data validation.
func (a *action[S]) validate(data *types.ActionData) (bool, error)

// Execute the action with receied data. Return the result of execution
// and error message, if any. The execution is only supposed to fail if
// and only if an error message is returned. The first response arg is
// just a status of the execution, so its meaning depends on contenxt.
func (a *action[S]) execute(data *types.ActionData) (bool, error)
