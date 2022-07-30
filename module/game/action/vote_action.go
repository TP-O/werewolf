package action

import (
	"uwwolf/module/game/core"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const VoteActionName = "Vote"

type voteAction struct {
	action[state.Poll]
}

func NewVoteAction(game core.Game) Action {
	voteAction := voteAction{
		action: action[state.Poll]{
			name:  VoteActionName,
			state: state.NewPoll([]types.PlayerId{1}),
			game:  game,
		},
	}

	return &voteAction
}

func (v *voteAction) validate(data *types.ActionData) (bool, error) {
	return true, nil
}

func (v *voteAction) execute(data *types.ActionData) (bool, error) {
	v.state.Vote(types.PlayerId(data.Actor), types.PlayerId(data.Targets[0]), 1)

	return true, nil
}
