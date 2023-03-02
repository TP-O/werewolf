package contract

import (
	"uwwolf/game/types"
)

// Game is game instace.
type Game interface {
	// ID returns game's ID.
	ID() types.GameID

	StatusID() types.GameStatusID

	// Poll returns the in-game poll management state.
	// Each specific faction has different poll to interact with.
	Poll(factionID types.FactionID) Poll

	// Scheduler returns turn manager.
	Scheduler() Scheduler

	// Player returns the player by given player ID.
	Player(playerID types.PlayerID) Player

	// PlayerIDsWithRoleID returns the player ID list has the
	// givent role ID.
	PlayerIDsWithRoleID(roleID types.RoleID) []types.PlayerID

	// PlayerIDsWithFactionID returns the player ID list has the given faction ID.
	PlayerIDsWithFactionID(factionID types.FactionID, onlyAlive bool) []types.PlayerID

	// PlayerIDsWithoutFactionID returns the player ID list doesn't have the given faction ID.
	PlayerIDsWithoutFactionID(factionID types.FactionID, onlyAlive bool) []types.PlayerID

	// Prepare sets up the game and returns completion time in milisecond.
	Prepare() int64

	// Start starts the game.
	Start() bool

	// Finish fishes the game.
	Finish() bool

	// Play activates the player's ability.
	Play(playerID types.PlayerID, req types.ActivateAbilityRequest) types.ActionResponse

	// KillPlayer kills the player by the given player ID.
	KillPlayer(playerID types.PlayerID, isExited bool) Player
}
