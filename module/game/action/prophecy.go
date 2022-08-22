package action

import (
	"github.com/go-playground/validator/v10"

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
	return p.action.overridePerform(p, req)
}

func (p *prophecy) Validate(req *types.ActionRequest) validator.ValidationErrorsTranslations {
	isIdentified := p.state.Identify(req.Targets[0]) != types.UnknownFaction

	if isIdentified {
		return map[string]string{
			types.AlertErrorField: "Already known identity!",
		}
	}

	return nil
}

// Check if a player is werewolf or not
func (p *prophecy) Execute(req *types.ActionRequest) *types.ActionResponse {
	factionId := p.game.GetPlayer(req.Targets[0]).GetFactionId()

	if factionId == types.WerewolfFaction {
		p.state.Acquire(req.Targets[0], types.WerewolfFaction)

		return &types.ActionResponse{
			Ok:   true,
			Data: true,
		}
	}

	p.state.Acquire(req.Targets[0], types.VillagerFaction)

	return &types.ActionResponse{
		Ok:   true,
		Data: false,
	}
}
