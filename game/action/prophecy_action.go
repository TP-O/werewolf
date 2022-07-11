package action

import (
	"fmt"

	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
)

func NewProphecyAction() itf.IAction {
	return &action{
		name: "Prophecy",
		kit: actionKit{
			validate: validateProphecy,
			execute:  executeProphecy,
			skip:     skipProphecy,
		},
	}
}

func validateProphecy(instruction *typ.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeProphecy(instruction *typ.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " pophesied " + instruction.Targets[0] + " is werewolf")

	return true
}

func skipProphecy(instruction *typ.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " skipped")

	return true
}
