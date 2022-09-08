package action

import (
	"uwwolf/module/game/contract"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const ProphecyActionName = "Prophecy"

type prophecy struct {
	action[state.Knowledge]
}

func NewProphecy(game contract.Game) contract.Action {
	prophecy := prophecy{
		action: action[state.Knowledge]{
			name:  ProphecyActionName,
			state: state.NewKnowledge(),
			game:  game,
		},
	}

	return &prophecy
}

func (p *prophecy) Perform(req *types.ActionRequest) *types.ActionResponse {
	return p.action.perform(p.validate, p.execute, req)
}

func (p *prophecy) validate(req *types.ActionRequest) (alert string) {
	if req.ActorId == req.TargetIds[0] {
		alert = "WTF! You don't know who are you?"
	} else if p.state.Identify(req.TargetIds[0]) != types.UnknownFaction {
		alert = "Already known identity!"
	}

	return
}

func (p *prophecy) execute(req *types.ActionRequest) *types.ActionResponse {
	factionId := p.game.Player(req.TargetIds[0]).FactionId()

	// Check if a player is werewolf or not
	if factionId == types.WerewolfFaction {
		p.state.Acquire(req.TargetIds[0], types.WerewolfFaction)

		return &types.ActionResponse{
			Ok:   true,
			Data: true,
		}
	}

	p.state.Acquire(req.TargetIds[0], types.VillagerFaction)

	return &types.ActionResponse{
		Ok:   true,
		Data: false,
	}
}
