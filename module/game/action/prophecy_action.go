package action

import (
	"errors"
	"uwwolf/module/game"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const ProphecyAction = "Prophecy"

type prophecyAction struct {
	action[state.Prediction]
}

func NewProphecyAction(game game.Game) Action[state.Prediction] {
	prophecyAction := prophecyAction{
		action: action[state.Prediction]{
			name:  ProphecyAction,
			state: state.NewPrediction(),
			game:  game,
		},
	}

	prophecyAction.action.receiveKit(&prophecyAction)

	return &prophecyAction
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

func (p *prophecyAction) skip(data *types.ActionData) (bool, error) {
	return true, nil
}
