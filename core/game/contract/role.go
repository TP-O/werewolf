package contract

import "uwwolf/game/types"

// Role represents a specific role in a game.
type Role interface {
	// ID returns role's ID.
	ID() types.RoleID

	// PhaseID returns role's active phase ID.
	// PhaseID() types.PhaseID

	// FactionID returns role's faction ID.
	FactionID() types.FactionID

	// TurnID returns role's turn order in active phase.
	// TurnID() types.TurnID

	// BeginRoundID returns round in which this role be able to
	// use its abilities.
	// BeginRoundID() types.RoundID

	// ActiveLimit returns remaining times this role can use the specific ability.
	// Returns total limit if the `index` is -1.
	ActiveLimit(index int) types.Limit

	// RegisterTurn adds role's turn to the game schedule.
	RegisterTurn()

	// BeforeDeath is triggered before killing this role.
	// If returns false, the player assigned it is saved.
	BeforeDeath() bool

	// AfterDeath is triggered after killing this role.
	AfterDeath()

	// ActivateAbility executes the action corresponding to the required ability.
	ActivateAbility(req *types.ActivateAbilityRequest) *types.ActionResponse
}
