package contract

import (
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type Game interface {
	// ID returns game's ID.
	ID() enum.GameID

	// Poll returns the in-game poll management state.
	// Each specific faction has different poll to interact with.
	Poll(factionID enum.FactionID) Poll

	// Scheduler returns turn manager.
	Scheduler() Scheduler

	// Player returns the player by given player ID.
	Player(playerID enum.PlayerID) Player

	// PlayerIDsByRoleID returns an array of player IDs by the
	// given role ID.
	PlayerIDsByRoleID(roleID enum.RoleID) []enum.PlayerID

	// PlayerIDsByFactionID returns an array of player IDs by the
	// given faction ID.
	PlayerIDsByFactionID(factionID enum.FactionID) []enum.PlayerID

	// WerewolfPlayerIDs returns an array of werewolf player IDs.
	WerewolfPlayerIDs() []enum.PlayerID

	// NonWerewolfPlayerIDs returns an array of non-werewolf player IDs.
	NonWerewolfPlayerIDs() []enum.PlayerID

	// AlivePlayerIDs returns an array of alive player IDs.
	AlivePlayerIDs(roleID enum.RoleID) []enum.PlayerID

	// Start starts the game and returns starting time in milisecond.
	Start() int64

	// Finish finishes the game.
	Finish() bool

	// UsePlayerRole uses player's role if the current turn is its.
	UsePlayerRole(playerID enum.PlayerID, req *types.UseRoleRequest) *types.ActionResponse

	// ConnectPlayer connects or disconnects the player from the game.
	ConnectPlayer(playerID enum.PlayerID, isConnected bool) bool

	// ExitPlayer exited the player from the game and marks as dead.
	ExitPlayer(playerID enum.PlayerID) bool

	// KillPlayer kills the player by given player ID.
	KillPlayer(playerID enum.PlayerID, isExited bool) Player
}
