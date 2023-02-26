package contract

import (
	"uwwolf/game/types"
)

type Player interface {
	// ID returns player's ID.
	ID() types.PlayerID

	// MainRoleID returns player's main role id.
	MainRoleID() types.RoleID

	// RoleIds returns player's assigned role ids.
	RoleIDs() []types.RoleID

	// Roles returns player's assigned roles.
	Roles() map[types.RoleID]Role

	// FactionID returns player's faction ID.
	FactionID() types.FactionID

	// IsDead checks if player is dead.
	IsDead() bool

	// SetFactionID assigns the player to the new faction.
	SetFactionID(factionID types.FactionID)

	// Die marks the player as dead and triggers roles's events.
	// The player will not be saved by them if `isExited` is true.
	Die(isExited bool) bool

	// Revive marks the player as undead.
	Revive() bool

	// AssignRole assigns the role to the player, and the faction can
	// be updated based on assigned role.
	AssignRole(roleID types.RoleID) (bool, error)

	// RevokeRole removes the role from the player, and the faction can
	// be updated based on removed role.
	RevokeRole(roleID types.RoleID) (bool, error)

	// ExecuteAction executes one of player's available ability.
	// The executed ability is selected based on the requested
	// action.
	ExecuteAction(req types.ExecuteActionRequest) types.ActionResponse
}
