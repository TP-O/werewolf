package game

import "uwwolf/game/contract"

// gamemaster controlls a game.
type gamemaster struct {
	game contract.Game
}

var _ contract.Gamemaster = (*gamemaster)(nil)

func NewGamemaster(game contract.Game) contract.Gamemaster {
	return &gamemaster{
		game,
	}
}

// notify

// winner decision

// record game
