package action

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const RecognitionActionName = "Recognition"

type recognition struct {
	action[[]types.PlayerId]

	knownRoleId types.RoleId
}

func NewRecognition(game contract.Game, roleId types.RoleId) contract.Action {
	recognition := recognition{
		action: action[[]types.PlayerId]{
			name:  RecognitionActionName,
			state: nil,
			game:  game,
		},
		knownRoleId: roleId,
	}

	return &recognition
}

func (r *recognition) Perform(req *types.ActionRequest) *types.ActionResponse {
	return r.action.perform(r.validate, r.execute, req)
}

func (r *recognition) validate(req *types.ActionRequest) (alert string) {
	if r.state != nil {
		alert = "Please double check your memory :D"
	}

	return
}

func (r *recognition) execute(req *types.ActionRequest) *types.ActionResponse {
	playerIds := r.game.PlayerIdsWithRole(r.knownRoleId)

	r.action.state = &playerIds

	return &types.ActionResponse{
		Ok:   true,
		Data: playerIds,
	}
}
