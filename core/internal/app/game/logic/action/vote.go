package action

import (
	"fmt"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type VoteActionSetting struct {
	FactionId types.FactionId
	PlayerId  types.PlayerId
	Weight    uint
}

type vote struct {
	action

	// poll is `Poll` instance on which this action votes.
	poll contract.Poll
}

func NewVote(world contract.World, setting *VoteActionSetting) (contract.Action, error) {
	if poll := world.Poll(setting.FactionId); poll == nil {
		return nil, fmt.Errorf("Poll does not exist ¯\\_(ツ)_/¯")
	} else {
		poll.AddElectors(setting.PlayerId)
		poll.SetWeight(setting.PlayerId, setting.Weight)

		return &vote{
			action: action{
				id:    VoteActionId,
				world: world,
			},
			poll: world.Poll(setting.FactionId),
		}, nil
	}
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (v *vote) Execute(req types.ActionRequest) types.ActionResponse {
	return v.action.execute(v, v.Id(), &req)
}

// skip ingores the action request.
func (v *vote) skip(req *types.ActionRequest) types.ActionResponse {
	// Abstain from voting
	if _, err := v.poll.Vote(req.ActorId, types.PlayerId("")); err != nil {
		return types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}
	return v.action.skip(req)
}

// perform completes the action request.
func (v *vote) perform(req *types.ActionRequest) types.ActionResponse {
	if ok, err := v.poll.Vote(req.ActorId, req.TargetId); !ok {
		return types.ActionResponse{
			Ok:      false,
			Message: err.Error(),
		}
	}

	return types.ActionResponse{
		Ok:   true,
		Data: req.TargetId,
	}
}
