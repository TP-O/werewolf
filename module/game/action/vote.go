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

func NewVote(game contract.Game, factionId types.FactionId, playerId types.PlayerId, weight uint) contract.Action {
	vote := vote{
		action: action[state.Poll]{
			name:  VoteActionName,
			state: game.Poll(factionId),
			game:  game,
		},
	}

	vote.state.SetWeight(playerId, weight)

	return &vote
}

func (v *vote) Perform(req *types.ActionRequest) *types.ActionResponse {
	return v.action.perform(v.validate, v.execute, req)
}

func (v *vote) validate(req *types.ActionRequest) (alert string) {
	if !req.IsSkipped {
		if !v.state.IsAllowed(req.ActorId) {
			alert = "Not allowed to vote :("
		} else if v.state.IsVoted(req.ActorId) {
			alert = "Already voted! Wait for next turn, OK?"
		}
	}

	return
}

func (v *vote) execute(req *types.ActionRequest) *types.ActionResponse {
	poorPlayerId := types.UnknownPlayer

	if !req.IsSkipped {
		poorPlayerId = req.TargetIds[0]
	}

	if !v.state.Vote(req.ActorId, poorPlayerId) {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   types.SystemErrorTag,
				Alert: "Poll isn't been opened yet!",
			},
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: poorPlayerId,
	}
}
