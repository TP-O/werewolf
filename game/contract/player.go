package contract

import (
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type Player interface {
	// ID returns player's ID.
	ID() enum.PlayerID

	// MainRoleID returns player's main role id.
	MainRoleID() enum.RoleID

	// RoleIds returns player's assigned role ids.
	RoleIDs() []enum.RoleID

	// Roles returns player's assigned roles.
	Roles() map[enum.RoleID]Role

	// FactionID returns player's faction ID.
	FactionID() enum.FactionID

	// IsDead checks if player is dead.
	IsDead() bool

	// SetFactionID assigns the player to the new faction.
	SetFactionID(factionID enum.FactionID)

	// Die marks the player as dead and triggers roles's events.
	// The player will not be saved by them if `isExited` is true.
	Die(isExited bool) bool

	// Revive marks the player as undead.
	Revive() bool

	// AssignRole assigns the role to the player, and the faction can
	// be updated based on assigned role.
	AssignRole(roleID enum.RoleID) (bool, error)

	// RevokeRole removes the role from the player, and the faction can
	// be updated based on removed role.
	RevokeRole(roleID enum.RoleID) (bool, error)

	// UseAbility executes one of player's available skills.
	// The executed skill is selected based on the requested
	// action.
	UseAbility(req *types.UseRoleRequest) *types.ActionResponse
}
