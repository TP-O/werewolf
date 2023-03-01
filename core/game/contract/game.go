package contract

import (
	"uwwolf/game/types"
)

// Game is game instace.
type Game interface {
	// ID returns game's ID.
	ID() types.GameID

	// Poll returns the in-game poll management state.
	// Each specific faction has different poll to interact with.
	Poll(factionID types.FactionID) Poll

	// Scheduler returns turn manager.
	Scheduler() Scheduler

	// Player returns the player by given player ID.
	Player(playerID types.PlayerID) Player

	// PlayerIDsByRoleID returns the player ID list by the
	// given role ID.
	PlayerIDsByRoleID(roleID types.RoleID) []types.PlayerID

	// PlayerIDsByFactionID returns the player ID list by the
	// given faction ID.
	PlayerIDsByFactionID(factionID types.FactionID) []types.PlayerID

	// WerewolfPlayerIDs returns the werewolf player ID list.
	WerewolfPlayerIDs() []types.PlayerID

	// NonWerewolfPlayerIDs returns the non-werewolf player ID list.
	NonWerewolfPlayerIDs() []types.PlayerID

	// AlivePlayerIDs returns the alive player ID list.
	// AlivePlayerIDs(roleID types.RoleID) []types.PlayerID

	// Start starts the game and returns the started time in milisecond.
	Start() int64

	// Finish fishes the game.
	Finish() bool

	// Play activates the player's ability.
	Play(playerID types.PlayerID, req types.ActivateAbilityRequest) types.ActionResponse

	// KillPlayer kills the player by the given player ID.
	KillPlayer(playerID types.PlayerID, isExited bool) Player
}
