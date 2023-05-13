package contract

import "uwwolf/internal/app/game/logic/types"

// Role represents a specific role in a game.
type Role interface {
	// ID returns role's ID.
	Id() types.RoleId

	// PhaseID returns role's active phase ID.
	// PhaseID() types.PhaseID

	// FactionID returns role's faction ID.
	FactionId() types.FactionId

	// ActiveTimes returns remaining times this role can use the specific ability.
	// Returns total limit if the `index` is -1.
	ActiveTimes(index int) types.Times

	// OnAssign is triggered when the role is assigned to a player.
	OnAfterAssign()

	// OnRevoke is triggered when the role is removed from a player.
	OnAfterRevoke()

	// OnBeforeDeath is triggered before killing this role.
	// If returns false, the player assigned it is saved.
	OnBeforeDeath() bool

	// OnAfterDeath is triggered after killing this role.
	OnAfterDeath()

	// ActivateAbility executes the action corresponding to the required ability.
	Use(req types.RoleRequest) types.RoleResponse
}
