package action

import (
	"github.com/go-playground/validator/v10"

	"uwwolf/module/game/contract"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

const ShootingActionName = "Shooting"

type shooting struct {
	action[state.Shotgun]
}

func NewShooting(game contract.Game) contract.Action {
	shootingAction := shooting{
		action: action[state.Shotgun]{
			name:  ShootingActionName,
			state: state.NewShotgun(),
			game:  game,
		},
	}

	return &shootingAction
}

func (s *shooting) Perform(req *types.ActionRequest) *types.ActionResponse {
	return s.action.overridePerform(s, req)
}

func (s *shooting) Validate(req *types.ActionRequest) validator.ValidationErrorsTranslations {
	if s.state.IsShot() {
		return map[string]string{
			types.AlertErrorField: "Already shoot!",
		}
	}

	return nil
}

func (s *shooting) Execute(req *types.ActionRequest) *types.ActionResponse {
	s.state.Shoot(req.Targets[0])

	return &types.ActionResponse{
		Ok: true,
	}
}
