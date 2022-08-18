package contract

import "uwwolf/types"

type Role interface {
	// Get roles' name.
	GetName() string

	// Do something after being voted. Return false if exonerated,
	// otherwise return true.
	AfterBeingVoted() bool

	// Do something before death
	AfterDeath()

	// Check condition is satisfied then if pass, activate skill
	// corresponding to this role based on game context.
	ActivateSkill(req *types.ActionRequest) *types.ActionResponse
}
