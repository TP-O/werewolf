package contract

import (
	"uwwolf/game/types"
)

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

	// PlayerIDsByRoleID returns an array of player IDs by the
	// given role ID.
	PlayerIDsByRoleID(roleID types.RoleID) []types.PlayerID

	// PlayerIDsByFactionID returns an array of player IDs by the
	// given faction ID.
	PlayerIDsByFactionID(factionID types.FactionID) []types.PlayerID

	// WerewolfPlayerIDs returns an array of werewolf player IDs.
	WerewolfPlayerIDs() []types.PlayerID

	// NonWerewolfPlayerIDs returns an array of non-werewolf player IDs.
	NonWerewolfPlayerIDs() []types.PlayerID

	// AlivePlayerIDs returns an array of alive player IDs.
	AlivePlayerIDs(roleID types.RoleID) []types.PlayerID

	// Start starts the game and returns starting time in milisecond.
	Start() int64

	// Finish finishes the game.
	Finish() bool

	// UsePlayerRole uses player's role if the current turn is its.
	UsePlayerRole(playerID types.PlayerID, req types.ExecuteActionRequest) types.ActionResponse

	// ConnectPlayer connects or disconnects the player from the game.
	ConnectPlayer(playerID types.PlayerID, isConnected bool) bool

	// ExitPlayer exited the player from the game and marks as dead.
	ExitPlayer(playerID types.PlayerID) bool

	// KillPlayer kills the player by given player ID.
	KillPlayer(playerID types.PlayerID) Player
}
