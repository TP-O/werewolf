package action

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
)

type vote struct {
	*action[contract.Poll]
}

func NewVote(game contract.Game, setting *types.VoteActionSetting) contract.Action {
	if poll := game.Poll(setting.FactionID); poll == nil ||
		!poll.AddElectors(setting.PlayerID) {

		return nil
	} else {
		poll.SetWeight(setting.PlayerID, setting.Weight)

		return &vote{
			&action[contract.Poll]{
				id:    config.VoteActionID,
				game:  game,
				state: poll,
			},
		}
	}
}

func (v *vote) Perform(req *types.ActionRequest) *types.ActionResponse {
	return v.action.perform(v.validate, v.execute, v.skip, req)
}

func (v *vote) validate(req *types.ActionRequest) (msg string) {
	if v.state == nil {
		msg = "You were rejected from the poll (╯‵□′)╯︵┻━┻"
	} else if !v.state.CanVote(req.ActorID) {
		msg = "Not allowed to vote :("
	}

	return
}

func (v *vote) skip(req *types.ActionRequest) *types.ActionResponse {
	v.state.Vote(req.ActorID, types.PlayerID(""))

	return v.action.skip(req)
}

func (v *vote) execute(req *types.ActionRequest) *types.ActionResponse {
	targetID := req.TargetIDs[0]

	if !v.state.Vote(req.ActorID, targetID) {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Unable to vote!",
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: targetID,
	}
}
