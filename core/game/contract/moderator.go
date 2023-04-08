package contract

import "uwwolf/game/types"

// Moderator controlls a game.
type Moderator interface {
	GameID() types.GameID

	// StartGame starts the game.
	StartGame() int64

	// FinishGame ends the game.
	FinishGame() bool

	MovePlayer(playerID types.PlayerID, x float64, y float64) (bool, error)

	// RequestPlay receives the play request from the player.
	RequestPlay(playerID types.PlayerID, req *types.ActivateAbilityRequest) *types.ActionResponse
}
