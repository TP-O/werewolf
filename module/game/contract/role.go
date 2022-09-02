package contract

import "uwwolf/types"

type Role interface {
	// Id returns role's id.
	Id() types.RoleId

	// FactionId returns id of faction to which role belongs.
	FactionId() types.FactionId

	// AfterBeingVoted is called after being voted. Return false
	// if exonerated, otherwise return true. Do something before
	// death.
	AfterBeingVoted() bool

	// AfterDeath is called after player dies. Do something after it.
	AfterDeath()

	// ActivateSkill checks condition is satisfied then performs
	// action corresponding to this role if everything is ok.
	// Each time a skill is successfully activated, its expiration
	// will be reduced by 1.
	ActivateSkill(req *types.ActionRequest) *types.ActionResponse
}
