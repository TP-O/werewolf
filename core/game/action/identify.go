package action

import (
	"fmt"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

// identify gets a player list with a specific role or faction.
type identify struct {
	action

	// isIdentified marks as identified or not.
	isIdentified bool

	// roleID is expected role ID.
	RoleID types.RoleID `json:"expected_role_id"`

	// factionID is expected faction ID.
	FactionID types.FactionID `json:"expected_faction_id"`

	// Role stores an array of player IDs having expected role ID.
	Role []types.PlayerID `json:"role_identification"`

	// Faction stores an array of player IDs having expected faction ID.
	Faction []types.PlayerID `json:"faction_identification"`
}

func NewRoleIdentify(game contract.Game, roleID types.RoleID) contract.Action {
	return &identify{
		action: action{
			id:   vars.IdentifyActionID,
			game: game,
		},
		RoleID: roleID,
		Role:   make([]types.PlayerID, 0),
	}
}

func NewFactionIdentify(game contract.Game, factionID types.FactionID) contract.Action {
	return &identify{
		action: action{
			id:   vars.IdentifyActionID,
			game: game,
		},
		FactionID: factionID,
		Faction:   make([]types.PlayerID, 0),
	}
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (i *identify) Execute(req *types.ActionRequest) *types.ActionResponse {
	return i.action.execute(i, req)
}

// validate checks if the action request is valid.
func (i identify) validate(req *types.ActionRequest) error {
	if i.isIdentified {
		return fmt.Errorf("You already recognized everyone ¯\\(º_o)/¯")
	}

	return nil
}

// perform completes the action request.
func (i *identify) perform(req *types.ActionRequest) *types.ActionResponse {
	i.isIdentified = true

	if !i.FactionID.IsUnknown() {
		i.Faction = i.game.AlivePlayerIDsWithFactionID(i.FactionID)
		return &types.ActionResponse{
			Ok:   true,
			Data: i.Faction,
		}
	} else if !i.RoleID.IsUnknown() {
		i.Role = i.game.AlivePlayerIDsWithRoleID(i.RoleID)
		return &types.ActionResponse{
			Ok:   true,
			Data: i.Role,
		}
	} else {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Unable to recognize!",
		}
	}
}
