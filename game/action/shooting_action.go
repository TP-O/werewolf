package action

import (
	"fmt"

	"uwwolf/game/contract"
)

func NewShootingAction() contract.Action {
	return &action{
		name: "Shooting",
		kit: actionKit{
			validate: validateShooting,
			execute:  executeShooting,
			skip:     skipShooting,
		},
	}
}

func validateShooting(instruction contract.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeShooting(instruction contract.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " shot " + instruction.Targets[0])

	return true
}

func skipShooting(instruction contract.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " skipped")

	return true
}
