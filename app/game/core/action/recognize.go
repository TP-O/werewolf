package action

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
)

type memory struct {
	Role    []types.PlayerID `json:"role"`
	Faction []types.PlayerID `json:"faction"`
}

type recognize struct {
	*action[*memory]
	isRecognized        bool
	recognizedRoleID    types.RoleID
	recognizedFactionID types.FactionID
}

func NewRoleRecognize(game contract.Game, roleID types.RoleID) contract.Action {
	return &recognize{
		action: &action[*memory]{
			id:   config.RecognizeActionID,
			game: game,
			state: &memory{
				Role: make([]types.PlayerID, 0),
			},
		},
		recognizedRoleID: roleID,
	}
}

func NewFactionRecognize(game contract.Game, factionID types.FactionID) contract.Action {
	return &recognize{
		action: &action[*memory]{
			id:   config.RecognizeActionID,
			game: game,
			state: &memory{
				Faction: make([]types.PlayerID, 0),
			},
		},
		recognizedFactionID: factionID,
	}
}

func (r *recognize) Perform(req *types.ActionRequest) *types.ActionResponse {
	return r.action.perform(r.validate, r.execute, r.skip, req)
}

func (r *recognize) execute(req *types.ActionRequest) *types.ActionResponse {
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

	return &types.ActionResponse{
		Ok:   true,
		Data: r.state,
	}
}
