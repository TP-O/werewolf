package action

import (
	"uwwolf/module/game/contract"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const VoteActionName = "Vote"

type vote struct {
	action[state.Poll]
}

func NewVote(game contract.Game, player contract.Player, weight uint) contract.Action {
	vote := vote{
		action: action[state.Poll]{
			name:  VoteActionName,
			state: game.Poll(player.FactionId()),
			game:  game,
		},
	}

	vote.state.SetWeight(player.Id(), weight)

	return &vote
}

func (v *vote) Perform(req *types.ActionRequest) *types.ActionResponse {
	return v.action.perform(v.validate, v.execute, req)
}

func (v *vote) validate(req *types.ActionRequest) (alert string) {
	if !req.IsSkipped {
		if !v.state.IsAllowed(req.Actor) {
			alert = "Already voted! Wait for next turn, OK?"
		}
	}

	return
}

func (v *vote) execute(req *types.ActionRequest) *types.ActionResponse {
	poorPlayerId := types.UnknownPlayer

	if !req.IsSkipped {
		poorPlayerId = req.Targets[0]
	}

	if !v.state.Vote(req.Actor, poorPlayerId) {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   types.SystemErrorTag,
				Alert: "Unknown error :(",
			},
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: poorPlayerId,
	}
}
