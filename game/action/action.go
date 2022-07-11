package action

import (
	"uwwolf/contract/typ"
	"uwwolf/validator"
)

type actionKit struct {
	validate func(*typ.ActionInstruction) bool
	execute  func(*typ.ActionInstruction) bool
	skip     func(*typ.ActionInstruction) bool
}

type action struct {
	name string
	kit  actionKit
}

func (a *action) GetName() string {
	return a.name
}

func (a *action) Perform(instruction *typ.ActionInstruction) bool {
	if !validator.ValidateStruct(instruction) || !a.kit.validate(instruction) {
		return false
	}

	if instruction.Skipped {
		return a.kit.skip(instruction)
	}

	return a.kit.execute(instruction)
}
