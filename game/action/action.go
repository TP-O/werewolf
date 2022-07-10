package action

import (
	"uwwolf/game/contract"
	"uwwolf/validator"
)

type actionKit struct {
	validate func(contract.ActionInstruction) bool
	execute  func(contract.ActionInstruction) bool
	skip     func(contract.ActionInstruction) bool
}

type action struct {
	name string
	kit  actionKit
}

func (a *action) GetName() string {
	return a.name
}

func (a *action) Perform(instruction contract.ActionInstruction) bool {
	if !validator.ValidateStruct(instruction) || !a.kit.validate(instruction) {
		return false
	}

	if instruction.Skipped {
		return a.kit.skip(instruction)
	}

	return a.kit.execute(instruction)
}
