package action

import (
	"fmt"

	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
)

func NewVoteAction() itf.IAction {
	return &action{
		name: "Vote",
		kit: actionKit{
			validate: validateVote,
			execute:  executeVote,
			skip:     skipVote,
		},
	}
}

func validateVote(instruction *typ.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeVote(instruction *typ.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " voted " + instruction.Targets[0])

	return true
}

func skipVote(instruction *typ.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " skipped")

	return true
}
