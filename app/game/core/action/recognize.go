package action

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
)

type recognize struct {
	*action[*types.RecognizeState]
	isRecognized        bool
	recognizedRoleID    types.RoleID
	recognizedFactionID types.FactionID
}

func NewRoleRecognize(game contract.Game, roleID types.RoleID) contract.Action {
	return &recognize{
		action: &action[*types.RecognizeState]{
			id:   config.RecognizeActionID,
			game: game,
			state: &types.RecognizeState{
				Role: make([]types.PlayerID, 0),
			},
		},
		recognizedRoleID: roleID,
	}
}

func NewFactionRecognize(game contract.Game, factionID types.FactionID) contract.Action {
	return &recognize{
		action: &action[*types.RecognizeState]{
			id:   config.RecognizeActionID,
			game: game,
			state: &types.RecognizeState{
				Faction: make([]types.PlayerID, 0),
			},
		},
		recognizedFactionID: factionID,
	}
}

func (r *recognize) Execute(req *types.ActionRequest) *types.ActionResponse {
	return r.action.combine(r.Skip, r.Validate, r.Perform, req)
}

func (r *recognize) Perform(req *types.ActionRequest) *types.ActionResponse {
	if !r.isRecognized {
		if !r.recognizedFactionID.IsUnknown() {
			r.state.Faction = r.game.PlayerIDsByFactionID(r.recognizedFactionID)
		} else if !r.recognizedRoleID.IsUnknown() {
			r.state.Role = r.game.PlayerIDsByRoleID(r.recognizedRoleID)
		} else {
			return &types.ActionResponse{
				Ok:      false,
				Message: "System error!",
			}
		}

		r.isRecognized = true
	}

	if !r.recognizedFactionID.IsUnknown() {
		return &types.ActionResponse{
			Ok:   true,
			Data: r.state.Faction,
		}
	}

	return &types.ActionResponse{
		Ok:   true,
		Data: r.state.Role,
	}
}
