package contract

import (
	"uwwolf/app/game/types"
)

type Game interface {
	// Id returns game id.
	ID() types.GameID

	// Poll returns the in-game poll management state.
	// Each faction has different poll to interact with.
	Poll(factionID types.FactionID) Poll

	Scheduler() Scheduler

	// Player returns the player by id.
	Player(playerID types.PlayerID) Player

	PlayerIDsByRoleID(roleID types.RoleID) []types.PlayerID

	// PlayersWithFaction return all players in a faction.
	PlayerIDsByFactionID(factionID types.FactionID) []types.PlayerID

	WerewolfPlayerIDs() []types.PlayerID

	NonWerewolfPlayerIDs() []types.PlayerID

	AlivePlayerIDs(roleID types.RoleID) []types.PlayerID

	// Start emits the signal to start game.
	Start() int64

	Finish() bool

	// RequestAction validates the request with its actor and then executes action
	// if everything is ok.
	UsePlayerRole(playerID types.PlayerID, req *types.UseRoleRequest) *types.ActionResponse

	ConnectPlayer(playerID types.PlayerID, isConnected bool) bool

	ExitPlayer(playerID types.PlayerID) bool

	KillPlayer(playerID types.PlayerID) Player
}
