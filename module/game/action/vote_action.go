package action

import (
	"uwwolf/module/game"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const VoteAction = "Vote"

type voteAction struct {
	action[state.Poll]
}

func NewVoteAction(game game.Game) Action[state.Poll] {
	voteAction := voteAction{
		action: action[state.Poll]{
			name:  ShootingAction,
			state: state.NewPoll([]types.PlayerId{1}),
			game:  game,
		},
	}

	voteAction.action.receiveKit(&voteAction)

	return &voteAction
}

func (v *voteAction) validate(data *types.ActionData) (bool, error) {
	return true, nil
}

func (v *voteAction) execute(data *types.ActionData) (bool, error) {
	v.state.Vote(types.PlayerId(data.Actor), types.PlayerId(data.Targets[0]), 1)

	return true, nil
}

func (v *voteAction) skip(data *types.ActionData) (bool, error) {
	v.state.Vote(types.PlayerId(data.Actor), 0, 0)

	return true, nil
}
