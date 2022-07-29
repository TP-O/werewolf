package action

import (
	"errors"
	"uwwolf/module/game"
	"uwwolf/types"
	"uwwolf/validator"
)

type action[S any] struct {
	name     string
	state    S
	game     game.Game
	validate func(*types.ActionData) (bool, error)
	execute  func(*types.ActionData) (bool, error)
	skip     func(*types.ActionData) (bool, error)
}

type Action[S any] interface {
	GetName() string
	GetState() S
	GetJson() string
	Perform(data *types.ActionData) (bool, error)
}

type actionKit interface {
	validate(*types.ActionData) (bool, error)
	execute(*types.ActionData) (bool, error)
	skip(*types.ActionData) (bool, error)
}

func (a *action[S]) GetName() string {
	return a.name
}

func (a *action[S]) GetState() S {
	return a.state
}

func (a *action[S]) GetJson() string {
	return ""
}

func (a *action[S]) Perform(data *types.ActionData) (bool, error) {
	if !validator.SimpleValidateStruct(data) {
		return false, errors.New("Invalid action!")
	}

	// Validate for each specific action
	if _, err := a.validate(data); err != nil {
		return false, err
	}

	if data.Skipped {
		return a.skip(data)
	}

	return a.execute(data)
}

// Assign 3 handle methods from embedding struct to this action,
// therefore; embedding struct can reuse method Perform of this action,
// but still keep pointer receiver to it.
func (a *action[S]) receiveKit(otherAction actionKit) {
	a.validate = otherAction.validate
	a.execute = otherAction.execute
	a.skip = otherAction.skip
}
