package action

import (
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
)

const ProphecyAction = "Prophecy"

func NewProphecyAction() itf.IAction {
	return &action[uint]{
		name: ProphecyAction,
		kit: actionKit[uint]{
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

func executeProphecy(_ itf.IGame, instruction *typ.ActionInstruction, _ uint) bool {
	// fmt.Println(instruction.Actor + " pophesied " + instruction.Targets[0] + " is werewolf")

	return true
}

func skipProphecy(_ itf.IGame, instruction *typ.ActionInstruction, _ uint) bool {
	// fmt.Println(instruction.Actor + " skipped")

	return true
}
