package contract

import "uwwolf/game/types"

// Moderator controlls a game.
type Moderator interface {
	// InitGame creates a new idle game instance.
	InitGame(newGame types.CreateGameRequest) bool

	// StartGame starts the owned game.
	StartGame() bool

	// FinishGame finishes the owned game.
	FinishGame() bool

	// ReceivePlayRequest receives ability request from the player.
	RequestPlay(playerID types.PlayerID, req *types.ActivateAbilityRequest) *types.ActionResponse
}
