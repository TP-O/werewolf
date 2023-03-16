package contract

import "uwwolf/game/types"

// Moderator controlls a game.
type Moderator interface {
	// InitGame creates a new idle game instance.
	InitGame(newGame *types.GameSetting) bool

	// StartGame starts the game.
	StartGame() bool

	// FinishGame ends the game.
	FinishGame() bool

	// RequestPlay receives the play request from the player.
	RequestPlay(playerID types.PlayerID, req *types.ActivateAbilityRequest) *types.ActionResponse
}
