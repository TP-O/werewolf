package action

import (
	"errors"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"
)

// identify gets a player list with a specific role or faction.
type identify struct {
	action

	// isIdentified marks as identified or not.
	isIdentified bool

	// roleId is expected role ID.
	RoleId types.RoleId `json:"expected_role_id"`

	// factionId is expected faction ID.
	FactionId types.FactionId `json:"expected_faction_id"`

	// Role stores an array of player IDs having expected role ID.
	Role []types.PlayerId `json:"role_identification"`

	// Faction stores an array of player IDs having expected faction ID.
	Faction []types.PlayerId `json:"faction_identification"`
}

func NewRoleIdentify(world contract.World, roleId types.RoleId) contract.Action {
	return &identify{
		action: action{
			id:    IdentifyActionId,
			world: world,
		},
		RoleId: roleId,
		Role:   make([]types.PlayerId, 0),
	}
}

func NewFactionIdentify(world contract.World, factionId types.FactionId) contract.Action {
	return &identify{
		action: action{
			id:    IdentifyActionId,
			world: world,
		},
		FactionId: factionId,
		Faction:   make([]types.PlayerId, 0),
	}
}

// Execute checks if the request is skipped. If so, skips the execution;
// otherwise, validates the request, and then performs the required action.
func (i *identify) Execute(req types.ActionRequest) types.ActionResponse {
	return i.action.execute(i, i.Id(), &req)
}

// validate checks if the action request is valid.
func (i identify) validate(req *types.ActionRequest) error {
	if i.isIdentified {
		return errors.New("You already recognized everyone ¯\\(º_o)/¯")
	}

	return nil
}

// perform completes the action request.
func (i *identify) perform(req *types.ActionRequest) types.ActionResponse {
	i.isIdentified = true

	if !util.IsZero(i.FactionId) {
		i.Faction = i.world.AlivePlayerIdsWithFactionId(i.FactionId)
		return types.ActionResponse{
			Ok:   true,
			Data: i.Faction,
		}
	} else if !util.IsZero(i.RoleId) {
		i.Role = i.world.AlivePlayerIdsWithRoleId(i.RoleId)
		return types.ActionResponse{
			Ok:   true,
			Data: i.Role,
		}
	} else {
		return types.ActionResponse{
			Ok:      false,
			Message: "Unable to recognize!",
		}
	}
}
