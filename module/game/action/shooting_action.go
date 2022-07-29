package action

import (
	"errors"
	"uwwolf/module/game"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const ShootingAction = "Shooting"

type shootingAction struct {
	action[state.Shotgun]
}

func NewShootingAction(game game.Game) Action[state.Shotgun] {
	shootingAction := shootingAction{
		action: action[state.Shotgun]{
			name:  ShootingAction,
			state: state.NewShotgun(),
			game:  game,
		},
	}

	shootingAction.action.receiveKit(&shootingAction)

	return &shootingAction
}

func (s *shootingAction) validate(data *types.ActionData) (bool, error) {
	if s.state.IsShot() {
		return false, errors.New("Already shoot!")
	}

	return true, nil
}

func (s *shootingAction) execute(data *types.ActionData) (bool, error) {
	// Check target in game

	s.state.Shoot(types.PlayerId(data.Targets[0]))

	return true, nil
}

func (s *shootingAction) skip(data *types.ActionData) (bool, error) {
	return true, nil
}
