package contract

import (
	"uwwolf/internal/app/game/logic/types"
)

// Game is game instace.
type World interface {
	// Poll returns the in-game poll management state.
	// Each specific faction has different poll to interact with.
	Poll(factionId types.FactionId) Poll

	// Scheduler returns turn manager.
	Scheduler() Scheduler

	// Player returns the player by given player ID.
	Player(playerID types.PlayerId) Player

	// Players returns the player list.
	Players() map[types.PlayerId]Player

	// AlivePlayerIDsWithRoleID returns the alive player ID list having the
	// givent role ID.
	AlivePlayerIdsWithRoleId(roleId types.RoleId) []types.PlayerId

	// AlivePlayerIDsWithFactionID returns the alive player ID list having the
	// given faction ID.
	AlivePlayerIdsWithFactionId(factionId types.FactionId) []types.PlayerId

	// AlivePlayerIDsWithoutFactionID returns the alive player ID list not having
	// the given faction ID.
	AlivePlayerIdsWithoutFactionId(factionId types.FactionId) []types.PlayerId

	// Load sets up the game and returns completion time in milisecond.
	Load() int64

	Map() Map
}

// type World interface {
// 	// Player returns player by ID.
// 	Player(ID types.PlayerId) types.PlayerEntity

// 	// AddObject adds the new entity to world.
// 	AddEntity(settings types.EntitySettings)

// 	// AddPlayer adds the new player entity to world.
// 	AddPlayer(ID types.PlayerId, x float64, y float64)

// 	// MovePlayer moves the player to specifjed location.
// 	// The new location had to be validated before moving.
// 	MovePlayer(ID types.PlayerId, x float64, y float64) (bool, error)
// }
