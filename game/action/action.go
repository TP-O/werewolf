package action

import (
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/validator"
)

type actionKit[T any] struct {
	validate func(*typ.ActionInstruction) bool
	execute  func(itf.IGame, *typ.ActionInstruction, *T) bool
	skip     func(itf.IGame, *typ.ActionInstruction, *T) bool
}

type action[T any] struct {
	name  string
	kit   actionKit[T]
	state T
}

func (a *action[T]) GetName() string {
	return a.name
}

func (a *action[T]) Perform(game itf.IGame, instruction *typ.ActionInstruction) bool {
	if !validator.ValidateStruct(instruction) || !a.kit.validate(instruction) {
		return false
	}

	if instruction.Skipped {
		return a.kit.skip(game, instruction, &a.state)
	}

	return a.kit.execute(game, instruction, &a.state)
}
