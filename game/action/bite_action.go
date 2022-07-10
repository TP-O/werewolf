package action

import (
	"fmt"

	"uwwolf/game/contract"
)

func NewBiteAction() contract.Action {
	return &action{
		name: "Bite",
		kit: actionKit{
			validate: validateBite,
			execute:  executeBite,
			skip:     skipBite,
		},
	}
}

func validateBite(instruction contract.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeBite(instruction contract.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " want to bite " + instruction.Targets[0])

	return true
}

func skipBite(instruction contract.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " skipped")

	return true
}
