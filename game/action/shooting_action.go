package action

import (
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
)

const ShootingAction = "Shooting"

func NewShootingAction() itf.IAction {
	return &action[uint]{
		name: "Shooting",
		kit: actionKit[uint]{
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

func executeShooting(_ itf.IGame, instruction *typ.ActionInstruction, _ uint) bool {
	// fmt.Println(instruction.Actor + " shot " + instruction.Targets[0])

	return true
}

func skipShooting(_ itf.IGame, instruction *typ.ActionInstruction, _ uint) bool {
	// fmt.Println(instruction.Actor + " skipped")

	return true
}
