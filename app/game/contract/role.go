package contract

import "uwwolf/app/types"

type Role interface {
	// Id returns role's id.
	Id() types.RoleId

	// PhaseId returns active phase of the role.
	PhaseId() types.PhaseId

	// FactionId returns id of faction to which role belongs.
	FactionId() types.FactionId

	// Priority returns role priority number in its phase.
	Priority() int

	// Score returns role score. The higher score role has, the easier
	// villager faction wins.
	Score() int

	// Set returns maximum number of this role can be appeared in game.
	// Returns -1 if it is infinite.
	Set() int

	// BeginRound return round id that the role can perform its action.
	BeginRound() types.RoundId

	// Expiration return number of times that this role can perform
	// its actions.
	Expiration() types.Expiration

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
