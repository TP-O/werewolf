package contract

import "uwwolf/types"

type Player interface {
	// Id returns player's id.
	Id() types.PlayerId

	// FactionId returns faction id's player.
	FactionId() types.FactionId

	// AssignRoles assigns a list of roles for the player,
	// and also updates FactionId based on assigned roles.
	AssignRoles(roles ...Role)

	// UseSkill executes one of player's available skills.
	// The executed skill is selected based on its settings.
	UseSkill(req *types.ActionRequest) *types.ActionResponse
}
