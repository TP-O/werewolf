package action

import (
	"fmt"

	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
)

func NewShootingAction() itf.IAction {
	return &action{
		name: "Shooting",
		kit: actionKit{
			validate: validateShooting,
			execute:  executeShooting,
			skip:     skipShooting,
		},
	}
}

func validateShooting(instruction *typ.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeShooting(instruction *typ.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " shot " + instruction.Targets[0])

	return true
}

func skipShooting(instruction *typ.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " skipped")

	return true
}
