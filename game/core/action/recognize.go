package action

import (
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type recognize struct {
	*action[*types.RecognizeState]
	isRecognized        bool
	recognizedRoleID    enum.RoleID
	recognizedFactionID enum.FactionID
}

func NewRoleRecognize(game contract.Game, roleID enum.RoleID) contract.Action {
	return &recognize{
		action: &action[*types.RecognizeState]{
			id:   enum.RecognizeActionID,
			game: game,
			state: &types.RecognizeState{
				Role: make([]enum.PlayerID, 0),
			},
		},
		recognizedRoleID: roleID,
	}
}

func NewFactionRecognize(game contract.Game, factionID enum.FactionID) contract.Action {
	return &recognize{
		action: &action[*types.RecognizeState]{
			id:   enum.RecognizeActionID,
			game: game,
			state: &types.RecognizeState{
				Faction: make([]enum.PlayerID, 0),
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
		if !enum.IsUnknownFactionID(r.recognizedFactionID) {
			r.state.Faction = r.game.PlayerIDsByFactionID(r.recognizedFactionID)
		} else if !enum.IsUnknownRoleID(r.recognizedRoleID) {
			r.state.Role = r.game.PlayerIDsByRoleID(r.recognizedRoleID)
		} else {
			return &types.ActionResponse{
				Ok:      false,
				Message: "System error!",
			}
		}

		r.isRecognized = true
	}

	if !enum.IsUnknownFactionID(r.recognizedFactionID) {
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
