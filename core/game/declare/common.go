package declare

import "uwwolf/game/types"

// Number of remaining times
const (
	UnlimitedTimes types.Times = iota - 1
	OutOfTimes
	Once
	Twice
)
