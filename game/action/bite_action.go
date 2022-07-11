package action

import (
	"fmt"

	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
)

func NewBiteAction() itf.IAction {
	return &action{
		name: "Bite",
		kit: actionKit{
			validate: validateBite,
			execute:  executeBite,
			skip:     skipBite,
		},
	}
}

func validateBite(instruction *typ.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeBite(instruction *typ.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " want to bite " + instruction.Targets[0])

	return true
}

func skipBite(instruction *typ.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " skipped")

	return true
}
