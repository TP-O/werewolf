package contract

import "uwwolf/internal/app/game/logic/types"

// Moderator controlls a game.
type Moderator interface {
	GameID() types.GameID

	// StatusID retusn current game status ID.
	GameStatus() types.GameStatusID

	World() World

	// StartGame starts the game.
	StartGame() int64

	// FinishGame ends the game.
	FinishGame() bool

	Player(ID types.PlayerId) Player

	Scheduler() Scheduler

	OnPhaseChanged(fn func(mod Moderator))

	RegisterActionExecution(regis types.ActionExecutionRegisteration)

	// RequestPlay receives the play request from the player.
	RequestPlay(playerID types.PlayerId, req *types.RoleRequest) *types.RoleResponse
}
