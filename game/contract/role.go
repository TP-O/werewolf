package contract

import (
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

// Role represents a specific role in the game.
type Role interface {
	// ID returns role's ID.
	ID() enum.RoleID

	// PhaseID returns role's active phase ID.
	PhaseID() enum.PhaseID

	// FactionID returns role's faction ID.
	FactionID() enum.FactionID

	// Priority returns role's priority in active phase.
	Priority() enum.Priority

	// BeginRound returns round in which the role be able to
	// use the abilities.
	BeginRound() enum.Round

	// ActiveLimit returns remaining times the role can use the abilities.
	// Returns total limit if the `actionID`` is 0.
	ActiveLimit(actionID enum.ActionID) enum.Limit

	// BeforeDeath does something special before killing the role.
	// Return false if saved, otherwise return true.
	BeforeDeath() bool

	// AfterDeath does something special after killing the role.
	AfterDeath()

	// UseAbility looks up the ability containing the required action and
	// then uses it.
	UseAbility(req *types.UseRoleRequest) *types.ActionResponse
}
