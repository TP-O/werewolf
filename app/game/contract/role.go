package contract

import "uwwolf/app/game/types"

type Role interface {
	// Id returns role's id.
	ID() types.RoleID

	// PhaseId returns active phase of the role.
	PhaseID() types.PhaseID

	// FactionId returns id of faction to which role belongs.
	FactionID() types.FactionID

	// Priority returns role priority number in its phase.
	Priority() types.Priority

	// BeginRound return round id that the role can perform its action.
	BeginRound() types.Round

	// Expiration return number of times that this role can perform
	// its actions.
	ActiveLimit() types.Limit

	// AfterBeingVoted is called after being voted. Return false
	// if exonerated, otherwise return true. Do something before
	// death.
	AfterBeingVoted() bool

	// AfterDeath is called after player dies. Do something after it.
	AfterDeath()

	// UseAbility checks condition is satisfied then performs
	// action corresponding to this role if everything is ok.
	// Each time a skill is successfully activated, its expiration
	// will be reduced by 1.
	UseAbility(req *types.UseRoleRequest) *types.ActionResponse
}
