package action

import (
	"errors"

	"uwwolf/module/game/core"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const ShootingActionName = "Shooting"

type shootingAction struct {
	action[state.Shotgun]
}

func NewShootingAction(game core.Game) Action {
	shootingAction := shootingAction{
		action: action[state.Shotgun]{
			name:  ShootingActionName,
			state: state.NewShotgun(),
			game:  game,
		},
	}

	return &shootingAction
}

func (s *shootingAction) validate(data *types.ActionData) (bool, error) {
	if s.state.IsShot() {
		return false, errors.New("Already shoot!")
	}

	return true, nil
}

func (s *shootingAction) execute(data *types.ActionData) (bool, error) {
	s.state.Shoot(types.PlayerId(data.Targets[0]))

	return true, nil
}
