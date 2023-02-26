package action

import (
	"fmt"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

type VoteActionSetting struct {
	FactionID types.FactionID
	PlayerID  types.PlayerID
	Weight    uint
}

type vote struct {
	action

	// poll is `Poll` instance on which this action votes.
	poll contract.Poll
}

func NewVote(game contract.Game, setting VoteActionSetting) (contract.Action, error) {
	if poll := game.Poll(setting.FactionID); poll == nil {
		return nil, fmt.Errorf("Poll does not exist ¯\\_(ツ)_/¯")
	} else if !poll.AddElectors(setting.PlayerID) {
		return nil, fmt.Errorf("Unable to join to the poll ಠ_ಠ")
	} else {
		poll.SetWeight(setting.PlayerID, setting.Weight)

		return &vote{
			action: action{
				id:   VoteActionID,
				game: game,
			},
			poll: game.Poll(setting.FactionID),
		}, nil
	}
}

func (v *vote) Execute(req types.ActionRequest) types.ActionResponse {
	return v.action.execute(v, req)
}

func (v *vote) skip(req types.ActionRequest) types.ActionResponse {
	// Abstain from voting
	v.poll.Vote(req.ActorID, types.PlayerID(""))
	return v.action.skip(req)
}

func (v *vote) perform(req types.ActionRequest) types.ActionResponse {
	if ok, err := v.poll.Vote(req.ActorID, req.TargetID); !ok {
		return types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}

	return types.ActionResponse{
		Ok: true,
		StateChanges: types.StateChanges{
			VotedPlayerID: req.TargetID,
		},
	}
}