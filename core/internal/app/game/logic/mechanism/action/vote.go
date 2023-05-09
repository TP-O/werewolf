package action

import (
	"fmt"
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/mechanism/contract"
	"uwwolf/internal/app/game/logic/tool"
	"uwwolf/internal/app/game/logic/types"
)

type VoteActionSetting struct {
	FactionID types.FactionID
	PlayerID  types.PlayerID
	Weight    uint
}

type vote struct {
	action

	// poll is `Poll` instance on which this action votes.
	poll tool.Poll
}

func NewVote(world contract.World, setting *VoteActionSetting) (contract.Action, error) {
	if poll := world.Poll(setting.FactionID); poll == nil {
		return nil, fmt.Errorf("Poll does not exist ¯\\_(ツ)_/¯")
	} else {
		poll.AddElectors(setting.PlayerID)
		poll.SetWeight(setting.PlayerID, setting.Weight)

		return &vote{
			action: action{
				id:    declare.VoteActionID,
				world: world,
			},
			poll: world.Poll(setting.FactionID),
		}, nil
	}
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (v *vote) Execute(req *types.ActionRequest) *types.ActionResponse {
	return v.action.execute(v, req)
}

// skip ingores the action request.
func (v *vote) skip(req *types.ActionRequest) *types.ActionResponse {
	// Abstain from voting
	if _, err := v.poll.Vote(req.ActorID, types.PlayerID("")); err != nil {
		return &types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}
	return v.action.skip(req)
}

// perform completes the action request.
func (v *vote) perform(req *types.ActionRequest) *types.ActionResponse {
	if ok, err := v.poll.Vote(req.ActorID, req.TargetID); !ok {
		return &types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: req.TargetID,
	}
}
