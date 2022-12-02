package action

import (
	"errors"
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
)

type vote struct {
	*action[contract.Poll]
}

func NewVote(game contract.Game, setting *types.VoteActionSetting) (contract.Action, error) {
	if poll := game.Poll(setting.FactionID); poll == nil {
		return nil, errors.New("Poll does not exist ¯\\_(ツ)_/¯")
	} else if !poll.AddElectors(setting.PlayerID) {
		return nil, errors.New("Unable to join to the poll ಠ_ಠ")
	} else {
		poll.SetWeight(setting.PlayerID, setting.Weight)

		return &vote{
			&action[contract.Poll]{
				id:    config.VoteActionID,
				game:  game,
				state: poll,
			},
		}, nil
	}
}

func (v *vote) Execute(req *types.ActionRequest) *types.ActionResponse {
	return v.action.combine(v.Skip, v.Validate, v.Perform, req)
}

func (v *vote) Validate(req *types.ActionRequest) error {
	if err := v.action.Validate(req); err != nil {
		return err
	}

	if !v.state.CanVote(req.ActorID) {
		return errors.New("Not allowed to vote ¯\\_(ツ)_/¯")
	}

	return nil
}

func (v *vote) Skip(req *types.ActionRequest) *types.ActionResponse {
	v.state.Vote(req.ActorID, types.PlayerID(""))

	return v.action.Skip(req)
}

func (v *vote) Perform(req *types.ActionRequest) *types.ActionResponse {
	targetID := req.TargetIDs[0]

	if !v.state.Vote(req.ActorID, targetID) {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Unable to vote (╯°□°)╯︵ ┻━┻",
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: targetID,
	}
}
