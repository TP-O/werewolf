package contract

import (
	"uwwolf/internal/app/game/logic/types"

	"github.com/paulmach/orb"
)

// Player represents the player in a game.
type Player interface {
	// Id returns player's ID.
	Id() types.PlayerId

	// MainRoleId returns player's main role ID.
	MainRoleId() types.RoleId

	// RoleIds returns player's assigned role IDs.
	RoleIds() []types.RoleId

	// Roles returns player's assigned roles.
	Roles() map[types.RoleId]Role

	// FactionId returns player's faction ID.
	FactionId() types.FactionId

	// IsDead checks if player is dead.
	IsDead() bool

	// Position returns the curernt location of the player.
	Location() (float64, float64)

	// SetFactionId assigns the player to the new faction.
	SetFactionId(factionId types.FactionId)

	// Die kills the player and triggers roles events.
	Die() bool

	// Exit kills the player and ignores any trigger preventing death.
	Exit() bool

	// AssignRole assigns role to the player, and the faction can
	// be updated based on the role.
	AssignRole(roleId types.RoleId) (bool, error)

	// RevokeRole removes the role from the player, and updates faction
	// if needed
	RevokeRole(roleId types.RoleId) (bool, error)

	// UseRole uses one of player's available ability.
	UseRole(req types.RoleRequest) types.RoleResponse

	// Move moves the player to the given location.
	Move(position orb.Point) (bool, error)
}
