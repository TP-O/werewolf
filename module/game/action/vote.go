package action

import (
	"github.com/go-playground/validator/v10"

	"uwwolf/module/game/contract"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const VoteActionName = "Vote"

type vote struct {
	action[state.Poll]
}

func NewVote(game contract.Game, poll *state.Poll, weight int) contract.Action {
	vote := vote{
		action: action[state.Poll]{
			name:  VoteActionName,
			state: poll,
			game:  game,
		},
	}

	return &vote
}

func (v *vote) Perform(req *types.ActionRequest) *types.ActionResponse {
	return v.action.overridePerform(v, req)
}

func (v *vote) Validate(req *types.ActionRequest) validator.ValidationErrorsTranslations {
	if v.state.IsVoted(req.Actor) {
		return map[string]string{
			types.AlertErrorField: "Already voted! Wait for next turn, OK?",
		}
	}

	return nil
}

func (v *vote) Execute(req *types.ActionRequest) *types.ActionResponse {
	v.state.Vote(req.Actor, req.Targets[0])

	return &types.ActionResponse{
		Ok: true,
	}
}
