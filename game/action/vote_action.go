package action

import (
	"fmt"

	"uwwolf/game/contract"
)

func NewVoteAction() contract.Action {
	return &action{
		name: "Vote",
		kit: actionKit{
			validate: validateVote,
			execute:  executeVote,
			skip:     skipVote,
		},
	}
}

func validateVote(instruction contract.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeVote(instruction contract.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " voted " + instruction.Targets[0])

	return true
}

func skipVote(instruction contract.ActionInstruction) bool {
	fmt.Println(instruction.Actor + " skipped")

	return true
}
