package contract

import (
	"uwwolf/game/types"
)

// Game is game instace.
type Game interface {
	// StatusID retusn current game status ID.
	StatusID() types.GameStatusID

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

	// Prepare sets up the game and returns completion time in milisecond.
	Prepare() int64

	// Start starts the game.
	Start() bool

	// Finish fishes the game.
	Finish() bool

	// Play activates the player's ability.
	Play(playerID types.PlayerID, req *types.ActivateAbilityRequest) *types.ActionResponse

	// KillPlayer kills the player by the given player ID.
	// If `isExited` is true, any trigger preventing death is ignored.
	KillPlayer(playerID types.PlayerID, isExited bool) Player
}
