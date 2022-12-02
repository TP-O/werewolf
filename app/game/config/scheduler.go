package config

import "uwwolf/app/game/types"

const (
	LastRound types.Round = iota - 1
	_
	FirstRound
)

const (
	NextPosition types.Position = iota - 3
	SortedPosition
	LastPosition
	FirstPosition
)
