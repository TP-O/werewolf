package action

import (
	"time"

	"uwwolf/module/game"
	"uwwolf/module/game/stuff"
	"uwwolf/types"
)

const VoteAction = "Vote"

func NewVoteAction(game game.IGame, factionId int, timeout time.Duration) Action {
	poll := &stuff.Poll{}
	poll.Init(game, factionId, timeout)

	return &action[*stuff.Poll]{
		name:     VoteAction,
		state:    poll,
		validate: validateVote,
		execute:  executeVote,
		skip:     skipVote,
	}
}

func validateVote(data *types.ActionData) bool {
	return data.Skipped ||
		(!data.Skipped && len(data.Targets) == 1)
}

func executeVote(game game.IGame, data *types.ActionData, poll *stuff.Poll) bool {
	poll.Start()
	poll.Vote(data.Actor, data.Targets[0])

	return true
}

func skipVote(game game.IGame, data *types.ActionData, poll *stuff.Poll) bool {
	poll.Start()

	return true
}
