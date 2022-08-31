package contract

import (
	"uwwolf/module/game/state"
	"uwwolf/types"
)

type Game interface {
	// Round returns the in-game turn management state.
	Round() *state.Round

	// Poll returns the in-game poll management state.
	// Each faction has different poll to interact with.
	Poll(factionId types.FactionId) *state.Poll

	// IsStarted returns true if game is started, false otherwise.
	IsStarted() bool

	// Start starts game by initializing all needed states and assigning
	// role for all players.
	Start() bool

	// Player returns the player by id.
	Player(playerId types.PlayerId) Player

	// PlayersWithRole return all players in a role.
	PlayersWithRole(roleId types.RoleId) []Player

	// KillPlayer marks the player as died then does something based on
	// he/she roles, if any.
	KillPlayer(playerId types.PlayerId) Player

	// RequestAction validates the request with its actor then executes action
	// if everything is ok.
	RequestAction(playerId types.PlayerId, req *types.ActionRequest) *types.ActionResponse
}
