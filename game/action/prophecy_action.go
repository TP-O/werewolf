package action

import (
	"fmt"

	"uwwolf/game/contract"
)

func NewProphecyAction() contract.Action {
	return &action{
		name: "Prophecy",
		kit: actionKit{
			validate: validateProphecy,
			execute:  executeProphecy,
			skip:     skipProphecy,
		},
	}
}

func validateProphecy(instruction contract.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeProphecy(instruction contract.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " pophesied " + instruction.Targets[0] + " is werewolf")

	return true
}

func skipProphecy(instruction contract.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " skipped")

	return true
}
