package action

import (
	"uwwolf/app/enum"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type vote struct {
	action[contract.Poll]
}

func NewVote(game contract.Game, setting *types.VoteActionSetting) contract.Action {
	vote := vote{
		action: action[contract.Poll]{
			id:         enum.VoteActionId,
			state:      game.Poll(setting.FactionId),
			game:       game,
			expiration: enum.UnlimitedTimes,
		},
	}

	vote.state.SetWeight(setting.PlayerId, setting.Weight)

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
	var poorPlayerId types.PlayerId

	if !req.IsSkipped {
		poorPlayerId = req.TargetIds[0]
	}

	if !v.state.Vote(req.ActorId, poorPlayerId) {
		return &types.ActionResponse{
			Ok:           false,
			PerformError: "Poll is not opened yet!",
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: poorPlayerId,
	}
}
