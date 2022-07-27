package action

import (
	"uwwolf/module/game"
	"uwwolf/types"
	"uwwolf/validator"
)

type action[T any] struct {
	name     string
	state    T
	validate func(*types.ActionData) bool
	execute  func(game.IGame, *types.ActionData, T) bool
	skip     func(game.IGame, *types.ActionData, T) bool
}

type Action interface {
	GetName() string
	Perform(game game.IGame, instruction *types.ActionData) bool
}

func (a *action[T]) GetName() string {
	return a.name
}

func (a *action[T]) Perform(game game.IGame, data *types.ActionData) bool {
	if !validator.SimpleValidateStruct(data) || !a.validate(data) {
		return false
	}

	if data.Skipped {
		return a.skip(game, data, a.state)
	}

	return a.execute(game, data, a.state)
}
