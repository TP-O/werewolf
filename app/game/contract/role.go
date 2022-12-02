package contract

import "uwwolf/app/game/types"

// Role represents a specific role in the game.
type Role interface {
	// ID returns role's ID.
	ID() types.RoleID

	// PhaseID returns role's active phase ID.
	PhaseID() types.PhaseID

	// FactionID returns role's faction ID.
	FactionID() types.FactionID

	// Priority returns role's priority in active phase.
	Priority() types.Priority

	// BeginRound returns round in which the role be able to
	// use the abilities.
	BeginRound() types.Round

	// ActiveLimit returns remaining times the role can use the abilities.
	ActiveLimit(actionID types.ActionID) types.Limit

	// BeforeDeath does something special before killing the role.
	// Return false if saved, otherwise return true.
	BeforeDeath() bool

	// AfterDeath does something special after killing the role.
	AfterDeath()

	// UseAbility looks up the ability containing the required action and
	// then uses it.
	UseAbility(req *types.UseRoleRequest) *types.ActionResponse
}
