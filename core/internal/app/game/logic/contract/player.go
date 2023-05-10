package contract

import (
	"uwwolf/internal/app/game/logic/types"

	"github.com/paulmach/orb"
)

// Player represents the player in a game.
type Player interface {
	// ID returns player's ID.
	Id() types.PlayerId

	// MainRoleID returns player's main role id.
	MainRoleId() types.RoleId

	// RoleIDs returns player's assigned role ids.
	RoleIds() []types.RoleId

	// Roles returns player's assigned roles.
	Roles() map[types.RoleId]Role

	// FactionID returns player's faction ID.
	FactionId() types.FactionId

	// IsDead checks if player is dead.
	IsDead() bool

	// Position returns curernt location of the player.
	Location() (float64, float64)

	// SetFactionID assigns the player to the new faction.
	SetFactionId(factionId types.FactionId)

	// Die kills the player and triggers roles events.
	Die() bool

	// Exit kills the player and ignores any trigger preventing death.
	Exit() bool

	// AssignRole assigns the role to the player, and the faction can
	// be updated based on this role.
	AssignRole(roleId types.RoleId) (bool, error)

	// RevokeRole removes the role from the player, and the faction can
	// be updated based on removed role.
	RevokeRole(roleId types.RoleId) (bool, error)

	// ActivateAbility executes one of player's available ability.
	// The executed ability is selected based on the requested
	// action.
	ActivateAbility(req *types.ActivateAbilityRequest) *types.ActionResponse

	// Move moves the player to the specified location.
	Move(position orb.Point) (bool, error)
}
