package action

import (
	"errors"

	"uwwolf/module/game/core"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const ProphecyActionName = "Prophecy"

type prophecyAction struct {
	action[state.Prediction]
}

func NewProphecyAction(game core.Game) Action {
	prophecyAction := prophecyAction{
		action: action[state.Prediction]{
			name:  ProphecyActionName,
			state: state.NewPrediction(),
			game:  game,
		},
	}

	return &prophecyAction
}

func (p *prophecyAction) Perform(data *types.ActionData) (bool, error) {
	return p.action.overridePerform(p, data)
}

func (p *prophecyAction) validate(data *types.ActionData) (bool, error) {
	isPredicted := p.state.Identify(data.Targets[0]) != types.UnknownFaction

	if isPredicted {
		return false, errors.New("Already known identity!")
	}

	return true, nil
}

// Check if a player is werewolf or not
func (p *prophecyAction) execute(data *types.ActionData) (bool, error) {
	faction := 0

	if faction == int(types.UnknownFaction) {
		return false, errors.New("Unknow faction!")
	}

	if faction == int(types.WerewolfFaction) {
		p.state.Add(types.WerewolfFaction, data.Targets[0])

		return true, nil
	}

	p.state.Add(types.VillageFaction, data.Targets[0])

	return false, nil
}
