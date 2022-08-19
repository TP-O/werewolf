package contract

import "uwwolf/types"

type Player interface {
	GetId() types.PlayerId

	GetFactionId() types.FactionId

	AssignRoles(roles ...Role)

	// Decide which skill to use based on game context.
	UseSkill(req *types.ActionRequest) *types.ActionResponse
}
