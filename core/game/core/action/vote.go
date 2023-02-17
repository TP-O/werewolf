package action

import (
	"fmt"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type vote struct {
	*action[contract.Poll]
}

func NewVote(game contract.Game, setting *types.VoteActionSetting) (contract.Action, error) {
	if poll := game.Poll(setting.FactionID); poll == nil {
		return nil, fmt.Errorf("Poll does not exist ¯\\_(ツ)_/¯")
	} else if !poll.AddElectors(setting.PlayerID) {
		return nil, fmt.Errorf("Unable to join to the poll ಠ_ಠ")
	} else {
		poll.SetWeight(setting.PlayerID, setting.Weight)

		return &vote{
			&action[contract.Poll]{
				id:    enum.VoteActionID,
				game:  game,
				state: poll,
			},
		}, nil
	}
}

func (v *vote) Execute(req *types.ActionRequest) *types.ActionResponse {
	return v.action.combine(v.Skip, v.Validate, v.Perform, req)
}

func (v *vote) Skip(req *types.ActionRequest) *types.ActionResponse {
	v.state.Vote(req.ActorID, enum.PlayerID(""))

	return v.action.Skip(req)
}

func (v *vote) Perform(req *types.ActionRequest) *types.ActionResponse {
	targetID := req.TargetIDs[0]

	if ok, err := v.state.Vote(req.ActorID, targetID); !ok {
		return &types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: targetID,
	}
}
