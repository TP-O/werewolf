package action

import (
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
	return s.action.perform(s.validate, s.execute, req)
}

func (s *shooting) validate(req *types.ActionRequest) (alert string) {
	if !req.IsSkipped {
		if !req.IsSkipped && req.Actor == req.Targets[0] {
			alert = "Please don't commit suicide :("
		} else if s.state.IsShot() {
			alert = "Already shot!"
		}
	}

	return
}

func (s *shooting) execute(req *types.ActionRequest) *types.ActionResponse {
	if req.IsSkipped {
		return &types.ActionResponse{
			Ok: true,
		}
	}

	if !s.state.Shoot(req.Targets[0]) {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   types.SystemErrorTag,
				Alert: "Unknown error :(",
			},
		}
	}

	killedPlayer := s.game.KillPlayer(req.Targets[0])

	return &types.ActionResponse{
		Ok:   true,
		Data: killedPlayer.Id(),
	}
}
