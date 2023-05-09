package contract

import (
	"uwwolf/internal/app/game/logic/types"

	"github.com/paulmach/orb"
)

// Player represents the player in a game.
type Player interface {
	// ID returns player's ID.
	ID() types.PlayerID

	// MainRoleID returns player's main role id.
	MainRoleID() types.RoleID

	// RoleIDs returns player's assigned role ids.
	RoleIDs() []types.RoleID

	// Roles returns player's assigned roles.
	Roles() map[types.RoleID]Role

	// FactionID returns player's faction ID.
	FactionID() types.FactionID

	// IsDead checks if player is dead.
	IsDead() bool

	// Position returns curernt location of the player.
	Location() (float64, float64)

	// SetFactionID assigns this player to the new faction.
	SetFactionID(factionID types.FactionID)

	// Die marks this player as dead and triggers roles events.
	// If `isExited` is true, any trigger preventing death is ignored.
	Die(isExited bool) bool

	// AssignRole assigns the role to the player, and the faction can
	// be updated based on this role.
	AssignRole(roleID types.RoleID) (bool, error)

	// RevokeRole removes the role from the player, and the faction can
	// be updated based on removed role.
	RevokeRole(roleID types.RoleID) (bool, error)

	// ActivateAbility executes one of player's available ability.
	// The executed ability is selected based on the requested
	// action.
	ActivateAbility(req *types.ActivateAbilityRequest) *types.ActionResponse

	// Move moves the player to the specified location.
	Move(position orb.Point) (bool, error)
}
