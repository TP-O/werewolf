package contract

import (
	"uwwolf/app/game/types"
)

type Player interface {
	// Id returns player's id.
	ID() types.PlayerID

	// MainRoleId returns the main role id of the player.
	MainRoleID() types.RoleID

	// RoleIds returns assigned role ids slice.
	RoleIDs() []types.RoleID

	// Roles returns assigned roles map.
	Roles() map[types.RoleID]Role

	// FactionId returns faction id's player.
	FactionID() types.FactionID

	IsDead() bool

	Die()

	Revive()

	SetFactionID(factionID types.FactionID)

	// AssignRoles assigns a list of roles to the player,
	// and also updates FactionId based on assigned roles.
	AssignRoles(roleIDs ...types.RoleID)

	RevokeRoles(roleIDs ...types.RoleID) bool

	// UseSkill executes one of player's available skills.
	// The executed skill is selected based on its action
	// settings.
	UseAbility(req *types.UseRoleRequest) *types.ActionResponse
}
