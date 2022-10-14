package action

import (
	"uwwolf/app/enum"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type shooting struct {
	action[*types.PlayerId]
}

func NewShooting(game contract.Game) contract.Action {
	shootingAction := shooting{
		action: action[*types.PlayerId]{
			id:         enum.ShootingActionId,
			state:      nil,
			game:       game,
			expiration: enum.OneTimes,
		},
	}

	return &shootingAction
}

func (s *shooting) Perform(req *types.ActionRequest) *types.ActionResponse {
	return s.action.perform(s.validate, s.execute, req)
}

func (s *shooting) validate(req *types.ActionRequest) (alert string) {
	if !req.IsSkipped {
		if !req.IsSkipped && req.ActorId == req.TargetIds[0] {
			alert = "Please don't commit suicide :("
		} else if s.state != nil {
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

	killedPlayer := s.game.KillPlayer(req.TargetIds[0])

	return &types.ActionResponse{
		Ok:   true,
		Data: killedPlayer.Id(),
	}
}
