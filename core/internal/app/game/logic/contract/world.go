package contract

import (
	"uwwolf/internal/app/game/logic/types"
)

// Game is game instace.
type World interface {
	// Poll returns the in-game poll management state.
	// Each specific faction has different poll to interact with.
	Poll(factionID types.FactionID) Poll

	// Scheduler returns turn manager.
	Scheduler() Scheduler

	// Player returns the player by given player ID.
	Player(playerID types.PlayerID) Player

	// Players returns the player list.
	Players() map[types.PlayerID]Player

	// AlivePlayerIDsWithRoleID returns the alive player ID list having the
	// givent role ID.
	AlivePlayerIDsWithRoleID(roleID types.RoleID) []types.PlayerID

	// AlivePlayerIDsWithFactionID returns the alive player ID list having the
	// given faction ID.
	AlivePlayerIDsWithFactionID(factionID types.FactionID) []types.PlayerID

	// AlivePlayerIDsWithoutFactionID returns the alive player ID list not having
	// the given faction ID.
	AlivePlayerIDsWithoutFactionID(factionID types.FactionID) []types.PlayerID

	// Load sets up the game and returns completion time in milisecond.
	Load() int64

	Map() Map
}

// type World interface {
// 	// Player returns player by ID.
// 	Player(ID types.PlayerID) types.PlayerEntity

// 	// AddObject adds the new entity to world.
// 	AddEntity(settings types.EntitySettings)

// 	// AddPlayer adds the new player entity to world.
// 	AddPlayer(ID types.PlayerID, x float64, y float64)

// 	// MovePlayer moves the player to specifjed location.
// 	// The new location had to be validated before moving.
// 	MovePlayer(ID types.PlayerID, x float64, y float64) (bool, error)
// }
