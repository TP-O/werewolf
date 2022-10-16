package contract

import (
	"uwwolf/app/types"
)

type Game interface {
	// Id returns game id.
	Id() types.GameId

	// Round returns the in-game turn management state.
	Round() Round

	// Poll returns the in-game poll management state.
	// Each faction has different poll to interact with.
	Poll(factionId types.FactionId) Poll

	// Init create game state by initializing all needed states and
	// assigning role for all players.
	Init() (map[types.PlayerId]Player, error)

	// Start emits the signal to start game.
	Start()

	// Player returns the player by id.
	Player(playerId types.PlayerId) Player

	// PlayersWithRole return all players in a role.
	PlayerIdsWithRole(roleId types.RoleId) []types.PlayerId

	// PlayersWithFaction return all players in a faction.
	PlayerIdsWithFaction(factionId types.FactionId) []types.PlayerId

	// KillPlayer marks the player as died then does something based on
	// he/she roles, if any.
	KillPlayer(playerId types.PlayerId) Player

	// RequestAction validates the request with its actor and then executes action
	// if everything is ok.
	RequestAction(req *types.ActionRequest) *types.ActionResponse
}
