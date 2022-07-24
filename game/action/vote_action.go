package action

import (
	"time"
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/game/stuff"
)

const VoteAction = "Vote"

func NewVoteAction(game itf.IGame, timeout time.Duration) itf.IAction {
	poll := &stuff.Poll{}
	poll.Init(game, timeout)

	return &action[*stuff.Poll]{
		name:  VoteAction,
		state: poll,
		kit: actionKit[*stuff.Poll]{
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

func executeVote(game itf.IGame, instruction *typ.ActionInstruction, poll *stuff.Poll) bool {
	poll.Start()
	poll.Vote(instruction.Actor, instruction.Targets[0])

	return true
}

func skipVote(game itf.IGame, instruction *typ.ActionInstruction, poll *stuff.Poll) bool {
	poll.Start()

	return true
}
