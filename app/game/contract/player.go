package contract

import "uwwolf/app/types"

type Player interface {
	// Id returns player's id.
	Id() types.PlayerId

	// MainRoleId returns the main role id of the player.
	MainRoleId() types.RoleId

	// RoleIds returns assigned role ids slice.
	RoleIds() []types.RoleId

	// Roles returns assigned roles map.
	Roles() map[types.RoleId]Role

	// FactionId returns faction id's player.
	FactionId() types.FactionId

	// AssignRoles assigns a list of roles to the player,
	// and also updates FactionId based on assigned roles.
	AssignRoles(roles ...Role)

	// UseSkill executes one of player's available skills.
	// The executed skill is selected based on its action
	// settings.
	UseSkill(req *types.ActionRequest) *types.ActionResponse
}
